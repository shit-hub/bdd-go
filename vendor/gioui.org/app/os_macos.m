// SPDX-License-Identifier: Unlicense OR MIT

// +build darwin,!ios

@import AppKit;

#include "_cgo_export.h"

__attribute__ ((visibility ("hidden"))) CALayer *gio_layerFactory(void);

@interface GioAppDelegate : NSObject<NSApplicationDelegate>
@end

@interface GioWindowDelegate : NSObject<NSWindowDelegate>
@end

@implementation GioWindowDelegate
- (void)windowWillMiniaturize:(NSNotification *)notification {
	NSWindow *window = (NSWindow *)[notification object];
	gio_onHide((__bridge CFTypeRef)window.contentView);
}
- (void)windowDidDeminiaturize:(NSNotification *)notification {
	NSWindow *window = (NSWindow *)[notification object];
	gio_onShow((__bridge CFTypeRef)window.contentView);
}
- (void)windowWillEnterFullScreen:(NSNotification *)notification {
	NSWindow *window = (NSWindow *)[notification object];
	gio_onFullscreen((__bridge CFTypeRef)window.contentView);
}
- (void)windowWillExitFullScreen:(NSNotification *)notification {
	NSWindow *window = (NSWindow *)[notification object];
	gio_onWindowed((__bridge CFTypeRef)window.contentView);
}
- (void)windowDidChangeScreen:(NSNotification *)notification {
	NSWindow *window = (NSWindow *)[notification object];
	CGDirectDisplayID dispID = [[[window screen] deviceDescription][@"NSScreenNumber"] unsignedIntValue];
	CFTypeRef view = (__bridge CFTypeRef)window.contentView;
	gio_onChangeScreen(view, dispID);
}
- (void)windowDidBecomeKey:(NSNotification *)notification {
	NSWindow *window = (NSWindow *)[notification object];
	gio_onFocus((__bridge CFTypeRef)window.contentView, 1);
}
- (void)windowDidResignKey:(NSNotification *)notification {
	NSWindow *window = (NSWindow *)[notification object];
	gio_onFocus((__bridge CFTypeRef)window.contentView, 0);
}
- (void)windowWillClose:(NSNotification *)notification {
	NSWindow *window = (NSWindow *)[notification object];
	window.delegate = nil;
	gio_onClose((__bridge CFTypeRef)window.contentView);
}
@end

static void handleMouse(NSView *view, NSEvent *event, int typ, CGFloat dx, CGFloat dy) {
	NSPoint p = [view convertPoint:[event locationInWindow] fromView:nil];
	if (!event.hasPreciseScrollingDeltas) {
		// dx and dy are in rows and columns.
		dx *= 10;
		dy *= 10;
	}
	// Origin is in the lower left corner. Convert to upper left.
	CGFloat height = view.bounds.size.height;
	gio_onMouse((__bridge CFTypeRef)view, typ, [NSEvent pressedMouseButtons], p.x, height - p.y, dx, dy, [event timestamp], [event modifierFlags]);
}

@interface GioView : NSView <CALayerDelegate,NSTextInputClient>
@end

@implementation GioView
- (void)setFrameSize:(NSSize)newSize {
	[super setFrameSize:newSize];
	[self setNeedsDisplay:YES];
}
// drawRect is called when OpenGL is used, displayLayer otherwise.
// Don't know why.
- (void)drawRect:(NSRect)r {
	gio_onDraw((__bridge CFTypeRef)self);
}
- (void)displayLayer:(CALayer *)layer {
	layer.contentsScale = self.window.backingScaleFactor;
	gio_onDraw((__bridge CFTypeRef)self);
}
- (CALayer *)makeBackingLayer {
	CALayer *layer = gio_layerFactory();
	layer.delegate = self;
	return layer;
}
- (void)mouseDown:(NSEvent *)event {
	handleMouse(self, event, MOUSE_DOWN, 0, 0);
}
- (void)mouseUp:(NSEvent *)event {
	handleMouse(self, event, MOUSE_UP, 0, 0);
}
- (void)middleMouseDown:(NSEvent *)event {
	handleMouse(self, event, MOUSE_DOWN, 0, 0);
}
- (void)middletMouseUp:(NSEvent *)event {
	handleMouse(self, event, MOUSE_UP, 0, 0);
}
- (void)rightMouseDown:(NSEvent *)event {
	handleMouse(self, event, MOUSE_DOWN, 0, 0);
}
- (void)rightMouseUp:(NSEvent *)event {
	handleMouse(self, event, MOUSE_UP, 0, 0);
}
- (void)mouseMoved:(NSEvent *)event {
	handleMouse(self, event, MOUSE_MOVE, 0, 0);
}
- (void)mouseDragged:(NSEvent *)event {
	handleMouse(self, event, MOUSE_MOVE, 0, 0);
}
- (void)scrollWheel:(NSEvent *)event {
	CGFloat dx = -event.scrollingDeltaX;
	CGFloat dy = -event.scrollingDeltaY;
	handleMouse(self, event, MOUSE_SCROLL, dx, dy);
}
- (void)keyDown:(NSEvent *)event {
	NSString *keys = [event charactersIgnoringModifiers];
	gio_onKeys((__bridge CFTypeRef)self, (__bridge CFTypeRef)keys, [event timestamp], [event modifierFlags], true);
	[self interpretKeyEvents:[NSArray arrayWithObject:event]];
}
- (void)keyUp:(NSEvent *)event {
	NSString *keys = [event charactersIgnoringModifiers];
	gio_onKeys((__bridge CFTypeRef)self, (__bridge CFTypeRef)keys, [event timestamp], [event modifierFlags], false);
}
- (void)insertText:(id)string {
	gio_onText((__bridge CFTypeRef)self, (__bridge CFTypeRef)string);
}
- (void)doCommandBySelector:(SEL)sel {
	// Don't pass commands up the responder chain.
	// They will end up in a beep.
}

- (BOOL)hasMarkedText {
	int res = gio_hasMarkedText((__bridge CFTypeRef)self);
	return res ? YES : NO;
}
- (NSRange)markedRange {
	return gio_markedRange((__bridge CFTypeRef)self);
}
- (NSRange)selectedRange {
	return gio_selectedRange((__bridge CFTypeRef)self);
}
- (void)unmarkText {
	gio_unmarkText((__bridge CFTypeRef)self);
}
- (void)setMarkedText:(id)string
        selectedRange:(NSRange)selRange
     replacementRange:(NSRange)replaceRange {
	NSString *str;
	// string is either an NSAttributedString or an NSString.
	if ([string isKindOfClass:[NSAttributedString class]]) {
		str = [string string];
	} else {
		str = string;
	}
	gio_setMarkedText((__bridge CFTypeRef)self, (__bridge CFTypeRef)str, selRange, replaceRange);
}
- (NSArray<NSAttributedStringKey> *)validAttributesForMarkedText {
	return nil;
}
- (NSAttributedString *)attributedSubstringForProposedRange:(NSRange)range
                                                actualRange:(NSRangePointer)actualRange {
	NSString *str = CFBridgingRelease(gio_substringForProposedRange((__bridge CFTypeRef)self, range, actualRange));
	return [[NSAttributedString alloc] initWithString:str attributes:nil];
}
- (void)insertText:(id)string
  replacementRange:(NSRange)replaceRange {
	NSString *str;
	// string is either an NSAttributedString or an NSString.
	if ([string isKindOfClass:[NSAttributedString class]]) {
		str = [string string];
	} else {
		str = string;
	}
	gio_insertText((__bridge CFTypeRef)self, (__bridge CFTypeRef)str, replaceRange);
}
- (NSUInteger)characterIndexForPoint:(NSPoint)point {
	return NSNotFound;
}
- (NSRect)firstRectForCharacterRange:(NSRange)range
                         actualRange:(NSRangePointer)actualRange {
	return NSZeroRect;
}
@end

// Delegates are weakly referenced from their peers. Nothing
// else holds a strong reference to our window delegate, so
// keep a single global reference instead.
static GioWindowDelegate *globalWindowDel;

static CVReturn displayLinkCallback(CVDisplayLinkRef dl, const CVTimeStamp *inNow, const CVTimeStamp *inOutputTime, CVOptionFlags flagsIn, CVOptionFlags *flagsOut, void *displayLinkContext) {
	gio_onFrameCallback(dl);
	return kCVReturnSuccess;
}

CFTypeRef gio_createDisplayLink(void) {
	CVDisplayLinkRef dl;
	CVDisplayLinkCreateWithActiveCGDisplays(&dl);
	CVDisplayLinkSetOutputCallback(dl, displayLinkCallback, nil);
	return dl;
}

int gio_startDisplayLink(CFTypeRef dl) {
	return CVDisplayLinkStart((CVDisplayLinkRef)dl);
}

int gio_stopDisplayLink(CFTypeRef dl) {
	return CVDisplayLinkStop((CVDisplayLinkRef)dl);
}

void gio_releaseDisplayLink(CFTypeRef dl) {
	CVDisplayLinkRelease((CVDisplayLinkRef)dl);
}

void gio_setDisplayLinkDisplay(CFTypeRef dl, uint64_t did) {
	CVDisplayLinkSetCurrentCGDisplay((CVDisplayLinkRef)dl, (CGDirectDisplayID)did);
}

void gio_hideCursor() {
	@autoreleasepool {
		[NSCursor hide];
	}
}

void gio_showCursor() {
	@autoreleasepool {
		[NSCursor unhide];
	}
}

void gio_setCursor(NSUInteger curID) {
	@autoreleasepool {
		switch (curID) {
			case 1:
				[NSCursor.arrowCursor set];
				break;
			case 2:
				[NSCursor.IBeamCursor set];
				break;
			case 3:
				[NSCursor.pointingHandCursor set];
				break;
			case 4:
				[NSCursor.crosshairCursor set];
				break;
			case 5:
				[NSCursor.resizeLeftRightCursor set];
				break;
			case 6:
				[NSCursor.resizeUpDownCursor set];
				break;
			case 7:
				[NSCursor.openHandCursor set];
				break;
			default:
				[NSCursor.arrowCursor set];
				break;
		}
	}
}

CFTypeRef gio_createWindow(CFTypeRef viewRef, CGFloat width, CGFloat height, CGFloat minWidth, CGFloat minHeight, CGFloat maxWidth, CGFloat maxHeight) {
	@autoreleasepool {
		NSRect rect = NSMakeRect(0, 0, width, height);
		NSUInteger styleMask = NSTitledWindowMask |
			NSResizableWindowMask |
			NSMiniaturizableWindowMask |
			NSClosableWindowMask;

		NSWindow* window = [[NSWindow alloc] initWithContentRect:rect
													   styleMask:styleMask
														 backing:NSBackingStoreBuffered
														   defer:NO];
		if (minWidth > 0 || minHeight > 0) {
			window.contentMinSize = NSMakeSize(minWidth, minHeight);
		}
		if (maxWidth > 0 || maxHeight > 0) {
			window.contentMaxSize = NSMakeSize(maxWidth, maxHeight);
		}
		[window setAcceptsMouseMovedEvents:YES];
		NSView *view = (__bridge NSView *)viewRef;
		[window setContentView:view];
		[window makeFirstResponder:view];
		window.releasedWhenClosed = NO;
		window.delegate = globalWindowDel;
		return (__bridge_retained CFTypeRef)window;
	}
}

CFTypeRef gio_createView(void) {
	@autoreleasepool {
		NSRect frame = NSMakeRect(0, 0, 0, 0);
		GioView* view = [[GioView alloc] initWithFrame:frame];
		view.wantsLayer = YES;
		view.layerContentsRedrawPolicy = NSViewLayerContentsRedrawDuringViewResize;
		return CFBridgingRetain(view);
	}
}

@implementation GioAppDelegate
- (void)applicationDidFinishLaunching:(NSNotification *)aNotification {
	[NSApp setActivationPolicy:NSApplicationActivationPolicyRegular];
	[NSApp activateIgnoringOtherApps:YES];
	gio_onFinishLaunching();
}
- (void)applicationDidHide:(NSNotification *)aNotification {
	gio_onAppHide();
}
- (void)applicationWillUnhide:(NSNotification *)notification {
	gio_onAppShow();
}
@end

void gio_main() {
	@autoreleasepool {
		[NSApplication sharedApplication];
		GioAppDelegate *del = [[GioAppDelegate alloc] init];
		[NSApp setDelegate:del];

		NSMenuItem *mainMenu = [NSMenuItem new];

		NSMenu *menu = [NSMenu new];
		NSMenuItem *hideMenuItem = [[NSMenuItem alloc] initWithTitle:@"Hide"
															  action:@selector(hide:)
													   keyEquivalent:@"h"];
		[menu addItem:hideMenuItem];
		NSMenuItem *quitMenuItem = [[NSMenuItem alloc] initWithTitle:@"Quit"
															  action:@selector(terminate:)
													   keyEquivalent:@"q"];
		[menu addItem:quitMenuItem];
		[mainMenu setSubmenu:menu];
		NSMenu *menuBar = [NSMenu new];
		[menuBar addItem:mainMenu];
		[NSApp setMainMenu:menuBar];

		globalWindowDel = [[GioWindowDelegate alloc] init];

		[NSApp run];
	}
}
