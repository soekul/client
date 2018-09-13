//
//  ShareViewController.m
//  KeybaseShare
//
//  Created by Michael Maxim on 8/31/18.
//  Copyright © 2018 Keybase. All rights reserved.
//

#import "ShareViewController.h"
#import "keybase/keybase.h"
#import "Pusher.h"
#import <MobileCoreServices/MobileCoreServices.h>
#import <AVFoundation/AVFoundation.h>
#import "Fs.h"

#if TARGET_OS_SIMULATOR
const BOOL isSimulator = YES;
#else
const BOOL isSimulator = NO;
#endif


@interface ShareViewController ()
@property NSDictionary* convTarget; // the conversation we will be sharing into
@property BOOL hasInited; // whether or not init has succeeded yet
@end

@implementation ShareViewController

- (BOOL)isContentValid {
    return self.hasInited && self.convTarget != nil;
}

// presentationAnimationDidFinish is called after the screen has rendered, and is the recommended place for loading data.
- (void)presentationAnimationDidFinish {
  dispatch_async(dispatch_get_global_queue(DISPATCH_QUEUE_PRIORITY_DEFAULT, 0), ^{
    BOOL skipLogFile = NO;
    NSError* error = nil;
    NSDictionary* fsPaths = [[FsHelper alloc] setupFs:skipLogFile setupSharedHome:NO];
    KeybaseExtensionInit(fsPaths[@"home"], fsPaths[@"sharedHome"], fsPaths[@"logFile"], @"prod", isSimulator, &error);
    if (error != nil) {
      dispatch_async(dispatch_get_main_queue(), ^{
        // If Init failed, then let's throw up our error screen.
        NSLog(@"Failed to init: %@", error);
        InitFailedViewController* initFailed = [InitFailedViewController alloc];
        [initFailed setDelegate:self];
        [self pushConfigurationViewController:initFailed];
      });
      return;
    }
    [self setHasInited:YES]; // Init is complete, we can use this to take down spinner on convo choice row
   
    NSString* jsonSavedConv = KeybaseExtensionGetSavedConv(); // result is in JSON format
    if ([jsonSavedConv length] > 0) {
      NSData* data = [jsonSavedConv dataUsingEncoding:NSUTF8StringEncoding];
      NSDictionary* conv = [NSJSONSerialization JSONObjectWithData:data options: NSJSONReadingMutableContainers error: &error];
      if (!conv) {
        NSLog(@"failed to parse saved conv: %@", error);
      } else {
        // Success reading a saved convo, set it and reload the items to it.
        [self setConvTarget:conv];
      }
    }
    dispatch_async(dispatch_get_main_queue(), ^{
      [self validateContent];
      [self reloadConfigurationItems];
    });
  });
}

-(void)initFailedClosed {
  // Just bail out of the extension if init failed
  [self cancel];
}

- (NSItemProvider*)firstSatisfiesTypeIdentifierCond:(NSArray*)attachments cond:(BOOL (^)(NSItemProvider*))cond {
  for (NSItemProvider* a in attachments) {
    if (cond(a)) {
      return a;
    }
  }
  return nil;
}

- (NSMutableArray*)allSatisfiesTypeIdentifierCond:(NSArray*)attachments cond:(BOOL (^)(NSItemProvider*))cond {
  NSMutableArray* res = [NSMutableArray array];
  for (NSItemProvider* a in attachments) {
    if (cond(a)) {
      [res addObject:a];
    }
  }
  return res;
}

- (BOOL)isWebURL:(NSItemProvider*)item {
  // "file URLs" also have type "url", but we want to treat them as files, not text.
  return (BOOL)([item hasItemConformingToTypeIdentifier:@"public.url"] && ![item hasItemConformingToTypeIdentifier:@"public.file-url"]);
}

// getSendableAttachments will get a list of messages we want to send from the share attempt. The flow is as follows:
// - If there is a URL item, we take it and only it.
// - If there is a text item, we take it and only it.
// - If there are none of the above, collect all the images and videos.
// - If we still don't have anything, select only the first item and hope for the best.
- (NSArray*)getSendableAttachments {
  NSExtensionItem *input = self.extensionContext.inputItems.firstObject;
  NSArray* attachments = [input attachments];
  NSMutableArray* res = [NSMutableArray array];
  NSItemProvider* item = [self firstSatisfiesTypeIdentifierCond:attachments cond:^(NSItemProvider* a) {
    return [self isWebURL:a];
  }];
  if (item) {
   [res addObject:item];
  }
  if ([res count] == 0) {
    item = [self firstSatisfiesTypeIdentifierCond:attachments cond:^(NSItemProvider* a) {
      return (BOOL)([a hasItemConformingToTypeIdentifier:@"public.text"]);
    }];
    if (item) {
      [res addObject:item];
    }
  }
  if ([res count] == 0) {
    res = [self allSatisfiesTypeIdentifierCond:attachments cond:^(NSItemProvider* a) {
      return (BOOL)([a hasItemConformingToTypeIdentifier:@"public.image"] || [a hasItemConformingToTypeIdentifier:@"public.movie"]);
    }];
  }
  if([res count] == 0 && attachments.firstObject != nil) {
    [res addObject:attachments.firstObject];
  }
  return res;
}

// loadPreviewView solves the problem of running out of memory when rendering the image previews on URLs. In some
// apps, loading URLs from them will crash our extension because of memory constraints. Instead of showing the image
// preview, just paste the text into the compose box. Otherwise, just do the normal thing.
- (UIView*)loadPreviewView {
  NSArray* items = [self getSendableAttachments];
  if ([items count] == 0) {
    return [super loadPreviewView];
  }
  NSItemProvider* item = items[0];
  if ([self isWebURL:item]) {
    [item loadItemForTypeIdentifier:@"public.url" options:nil completionHandler:^(NSURL *url, NSError *error) {
      dispatch_async(dispatch_get_main_queue(), ^{
        [self.textView setText:[NSString stringWithFormat:@"%@\n%@", self.contentText, [url absoluteString]]];
      });
    }];
    return nil;
  }
  return [super loadPreviewView];
}

- (void)didReceiveMemoryWarning {
    KeybaseExtensionForceGC(); // run Go GC and hope for the best
    [super didReceiveMemoryWarning];
}

- (void)createVideoPreview:(NSURL*)url resultCb:(void (^)(int,int,int,int,int,NSData*))resultCb  {
  NSError *error = NULL;
  CMTime time = CMTimeMake(1, 1);
  AVURLAsset *asset = [[AVURLAsset alloc] initWithURL:url options:nil];
  AVAssetImageGenerator *generateImg = [[AVAssetImageGenerator alloc] initWithAsset:asset];
  [generateImg setAppliesPreferredTrackTransform:YES];
  CGImageRef cgOriginal = [generateImg copyCGImageAtTime:time actualTime:NULL error:&error];
  [generateImg setMaximumSize:CGSizeMake(640, 640)];
  CGImageRef cgThumb = [generateImg copyCGImageAtTime:time actualTime:NULL error:&error];
  int duration = CMTimeGetSeconds([asset duration]);
  UIImage* original = [UIImage imageWithCGImage:cgOriginal];
  UIImage* scaled = [UIImage imageWithCGImage:cgThumb];
  NSData* preview = UIImageJPEGRepresentation(scaled, 0.7);
  resultCb(duration, original.size.width, original.size.height, scaled.size.width, scaled.size.height, preview);
  CGImageRelease(cgOriginal);
  CGImageRelease(cgThumb);
}

- (void)createImagePreview:(NSURL*)url resultCb:(void (^)(int,int,int,int,NSData*))resultCb  {
  UIImage* original = [UIImage imageWithData:[NSData dataWithContentsOfURL:url]];
  CFURLRef cfurl = CFBridgingRetain(url);
  CGImageSourceRef is = CGImageSourceCreateWithURL(cfurl, nil);
  NSDictionary* opts = [[NSDictionary alloc] initWithObjectsAndKeys:
                        (id)kCFBooleanTrue, (id)kCGImageSourceCreateThumbnailWithTransform,
                        (id)kCFBooleanTrue, (id)kCGImageSourceCreateThumbnailFromImageIfAbsent,
                        [NSNumber numberWithInt:640], (id)kCGImageSourceThumbnailMaxPixelSize,
                        nil];
  CGImageRef image = CGImageSourceCreateThumbnailAtIndex(is, 0, (CFDictionaryRef)opts);
  UIImage* scaled = [UIImage imageWithCGImage:image];
  NSData* preview = UIImageJPEGRepresentation(scaled, 0.7);
  resultCb(original.size.width, original.size.height, scaled.size.width, scaled.size.height, preview);
  CGImageRelease(image);
  CFRelease(cfurl);
  CFRelease(is);
}

- (void) maybeCompleteRequest:(BOOL)lastItem {
  if (!lastItem) { return; }
  dispatch_async(dispatch_get_main_queue(), ^{
    [self.extensionContext completeRequestReturningItems:nil completionHandler:nil];
  });
}

// processItem will invokve the correct function on the Go side for the given attachment type.
- (void)processItem:(NSItemProvider*)item lastItem:(BOOL)lastItem {
  PushNotifier* pusher = [[PushNotifier alloc] init];
  NSString* convID = self.convTarget[@"ConvID"];
  NSString* name = self.convTarget[@"Name"];
  NSNumber* membersType = self.convTarget[@"MembersType"];
  NSItemProviderCompletionHandler urlHandler = ^(NSURL* url, NSError* error) {
    KeybaseExtensionPostText(convID, name, NO, [membersType longValue], self.contentText, pusher, &error);
    [self maybeCompleteRequest:lastItem];
  };
  
  NSItemProviderCompletionHandler textHandler = ^(NSString* text, NSError* error) {
    KeybaseExtensionPostText(convID, name, NO, [membersType longValue], text, pusher, &error);
    [self maybeCompleteRequest:lastItem];
  };
  
  // The NSItemProviderCompletionHandler interface is a little tricky. The caller of our handler
  // will inspect the arguments that we have given, and will attempt to give us the attachment
  // in this form. For files, we always want a file URL, and so that is what we pass in.
  NSItemProviderCompletionHandler fileHandler = ^(NSURL* url, NSError* error) {
    // Check for no URL (it might have not been possible for the OS to give us one)
    if (url == nil) {
      [self maybeCompleteRequest:lastItem];
      return;
    }
    NSString* filePath = [url relativePath];
    if ([item hasItemConformingToTypeIdentifier:@"public.movie"]) {
      // Generate image preview here, since it runs out of memory easy in Go
      [self createVideoPreview:url resultCb:^(int duration, int baseWidth, int baseHeight, int previewWidth, int previewHeight, NSData* preview) {
        NSError* error = NULL;
        KeybaseExtensionPostVideo(convID, name, NO, [membersType longValue], self.contentText, filePath,
                                 duration, baseWidth, baseHeight, previewWidth, previewHeight, preview, pusher, &error);
      }];
    } else if ([item hasItemConformingToTypeIdentifier:@"public.image"]) {
      // Generate image preview here, since it runs out of memory easy in Go
      [self createImagePreview:url resultCb:^(int baseWidth, int baseHeight, int previewWidth, int previewHeight, NSData* preview) {
        NSError* error = NULL;
        KeybaseExtensionPostJPEG(convID, name, NO, [membersType longValue], self.contentText, filePath,
                                 baseWidth, baseHeight, previewWidth, previewHeight, preview, pusher, &error);
      }];
    } else {
      NSError* error = NULL;
      KeybaseExtensionPostFile(convID, name, NO, [membersType longValue], self.contentText, filePath, pusher, &error);
    }
    [self maybeCompleteRequest:lastItem];
  };
  
  if ([item hasItemConformingToTypeIdentifier:@"public.movie"]) {
    [item loadItemForTypeIdentifier:@"public.movie" options:nil completionHandler:fileHandler];
  } else if ([item hasItemConformingToTypeIdentifier:@"public.image"]) {
    [item loadItemForTypeIdentifier:@"public.image" options:nil completionHandler:fileHandler];
  } else if ([item hasItemConformingToTypeIdentifier:@"public.file-url"]) {
    [item loadItemForTypeIdentifier:@"public.file-url" options:nil completionHandler:fileHandler];
  } else if ([item hasItemConformingToTypeIdentifier:@"public.text"]) {
    [item loadItemForTypeIdentifier:@"public.text" options:nil completionHandler:textHandler];
  } else if ([item hasItemConformingToTypeIdentifier:@"public.url"]) {
    [item loadItemForTypeIdentifier:@"public.url" options:nil completionHandler:urlHandler];
  } else {
    [pusher localNotification:@"extension" msg:@"We failed to send your message. Please try from the Keybase app."
                   badgeCount:-1 soundName:@"default" convID:@"" typ:@"chat.extension"];
    [self maybeCompleteRequest:lastItem];
  }
}

- (void)didSelectPost {
  if (!self.convTarget) {
    // Just bail out of here if nothing was selected
    [self maybeCompleteRequest:YES];
    return;
  }
  NSArray* items = [self getSendableAttachments];
  if ([items count] == 0) {
    [self maybeCompleteRequest:YES];
    return;
  }
  for (int i = 0; i < [items count]; i++) {
    BOOL lastItem = (BOOL)(i == [items count]-1);
    [self processItem:items[i] lastItem:lastItem];
  }
}

- (NSArray *)configurationItems {
  SLComposeSheetConfigurationItem *item = [[SLComposeSheetConfigurationItem alloc] init];
  item.title = @"Share to...";
  item.valuePending = !self.hasInited; // show a spinner if we haven't inited
  if (self.convTarget) {
    item.value = self.convTarget[@"Name"];
  } else if (self.hasInited) {
    item.value = @"Please choose";
  }
  item.tapHandler = ^{
    ConversationViewController *viewController = [[ConversationViewController alloc] initWithStyle:UITableViewStylePlain];
    viewController.delegate = self;
    [self pushConfigurationViewController:viewController];
  };
  return @[item];
}

- (void)convSelected:(NSDictionary *)conv {
  // This is a delegate method from the inbox view, it gets run when the user taps an item.
  [self setConvTarget:conv];
  [self validateContent];
  [self reloadConfigurationItems];
  [self popConfigurationViewController];
}

@end
