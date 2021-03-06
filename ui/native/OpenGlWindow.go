package native

import (
	"time"

	"github.com/go-gl/glfw/v3.2/glfw"

	"github.com/inkyblackness/hacked/ui/input"
	"github.com/inkyblackness/hacked/ui/opengl"
)

var buttonsByIndex = map[glfw.MouseButton]uint32{
	glfw.MouseButton1: input.MousePrimary,
	glfw.MouseButton2: input.MouseSecondary}

// OpenGLWindow represents a native OpenGL surface.
type OpenGLWindow struct {
	opengl.WindowEventDispatcher

	keyBuffer *input.StickyKeyBuffer

	glfwWindow *glfw.Window
	glWrapper  *OpenGL

	framesPerSecond float64
	frameTime       time.Duration
	nextRenderTick  time.Time
}

// NewOpenGLWindow tries to initialize the OpenGL environment and returns a
// new window instance.
func NewOpenGLWindow(title string, framesPerSecond float64) (window *OpenGLWindow, err error) {
	if err = glfw.Init(); err == nil {
		glfw.WindowHint(glfw.Resizable, 1)
		glfw.WindowHint(glfw.Decorated, 1)
		glfw.WindowHint(glfw.ClientAPI, glfw.OpenGLAPI)
		glfw.WindowHint(glfw.ContextVersionMajor, 3)
		glfw.WindowHint(glfw.ContextVersionMinor, 2)
		glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
		glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
		var glfwWindow *glfw.Window
		glfwWindow, err = glfw.CreateWindow(1280, 720, title, nil, nil)
		if err == nil {
			glfwWindow.MakeContextCurrent()

			window = &OpenGLWindow{
				WindowEventDispatcher: opengl.NullWindowEventDispatcher(),
				glfwWindow:            glfwWindow,
				glWrapper:             NewOpenGL(),
				framesPerSecond:       framesPerSecond,
				frameTime:             time.Duration(int64(float64(time.Second) / framesPerSecond)),
				nextRenderTick:        time.Now()}

			window.keyBuffer = input.NewStickyKeyBuffer(window.StickyKeyListener())

			glfwWindow.SetCursorPosCallback(window.onCursorPos)
			glfwWindow.SetMouseButtonCallback(window.onMouseButton)
			glfwWindow.SetScrollCallback(window.onMouseScroll)
			glfwWindow.SetFramebufferSizeCallback(window.onFramebufferResize)
			glfwWindow.SetKeyCallback(window.onKey)
			glfwWindow.SetCharCallback(window.onChar)
			glfwWindow.SetDropCallback(window.onDrop)
			glfwWindow.SetCloseCallback(window.onClosing)
		}
	}
	return
}

// ShouldClose returns true if the user requested the window to close.
func (window *OpenGLWindow) ShouldClose() bool {
	return window.glfwWindow.ShouldClose()
}

// Close closes the window and releases its resources.
func (window *OpenGLWindow) Close() {
	window.CallClosed()
	window.glfwWindow.Destroy()
	glfw.Terminate()
}

// ClipboardString returns the current value of the clipboard, if it is compatible with UTF-8.
func (window OpenGLWindow) ClipboardString() (string, error) {
	return window.glfwWindow.GetClipboardString()
}

// SetClipboardString sets the current value of the clipboard as UTF-8 string.
func (window OpenGLWindow) SetClipboardString(value string) {
	window.glfwWindow.SetClipboardString(value)
}

// Update must be called from within the main thread as often as possible.
func (window *OpenGLWindow) Update() {
	glfw.PollEvents()

	now := time.Now()
	delta := now.Sub(window.nextRenderTick)
	if delta.Nanoseconds() < 0 {
		// detected a change of wallclock time into the past; realign
		delta = window.frameTime
		window.nextRenderTick = now
	}

	if delta.Nanoseconds() >= window.frameTime.Nanoseconds() {
		window.glfwWindow.MakeContextCurrent()
		window.CallRender()
		window.glfwWindow.SwapBuffers()
		framesCovered := delta.Nanoseconds() / window.frameTime.Nanoseconds()
		window.nextRenderTick = window.nextRenderTick.Add(time.Duration(framesCovered * window.frameTime.Nanoseconds()))
	}
}

// OpenGL returns the OpenGL API.
func (window *OpenGLWindow) OpenGL() opengl.OpenGL {
	return window.glWrapper
}

// Size returns the dimension of the frame buffer of this window.
func (window *OpenGLWindow) Size() (width int, height int) {
	return window.glfwWindow.GetFramebufferSize()
}

// SetCursorVisible toggles the visibility of the cursor.
func (window *OpenGLWindow) SetCursorVisible(visible bool) {
	if visible {
		window.glfwWindow.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
	} else {
		window.glfwWindow.SetInputMode(glfw.CursorMode, glfw.CursorHidden)
	}
}

// SetFullScreen toggles the windowed mode.
func (window *OpenGLWindow) SetFullScreen(on bool) {
	if on {
		monitor := glfw.GetPrimaryMonitor()
		videoMode := monitor.GetVideoMode()
		window.glfwWindow.SetMonitor(monitor, 0, 0, videoMode.Width, videoMode.Height, glfw.DontCare)
	} else {
		window.glfwWindow.SetMonitor(nil, 0, 0, 1280, 720, glfw.DontCare)
	}
}

// SetCloseRequest sets the should-close property of the window.
func (window *OpenGLWindow) SetCloseRequest(shouldClose bool) {
	window.glfwWindow.SetShouldClose(shouldClose)
}

func (window *OpenGLWindow) onClosing(rawWindow *glfw.Window) {
	window.CallClosing()
}

func (window *OpenGLWindow) onFramebufferResize(rawWindow *glfw.Window, width int, height int) {
	window.CallResize(width, height)
}

func (window *OpenGLWindow) onCursorPos(rawWindow *glfw.Window, x float64, y float64) {
	window.CallOnMouseMove(float32(x), float32(y))
}

func (window *OpenGLWindow) onMouseButton(rawWindow *glfw.Window, rawButton glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	button, knownButton := buttonsByIndex[rawButton]

	if knownButton {
		modifier := window.mapModifier(mods)

		if action == glfw.Press {
			window.CallOnMouseButtonDown(button, modifier)
		} else if action == glfw.Release {
			window.CallOnMouseButtonUp(button, modifier)
		}
	}
}

func (window *OpenGLWindow) onMouseScroll(rawWindow *glfw.Window, dx float64, dy float64) {
	window.CallOnMouseScroll(float32(dx), float32(dy))
}

func (window *OpenGLWindow) onKey(rawWindow *glfw.Window, glfwKey glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	modifier := window.mapModifier(mods)
	key, knownKey := keyMap[glfwKey]

	if knownKey {
		switch action {
		case glfw.Press:
			window.keyBuffer.KeyDown(key, modifier)
		case glfw.Repeat:
			window.keyBuffer.KeyUp(key, modifier)
			window.keyBuffer.KeyDown(key, modifier)
		case glfw.Release:
			window.keyBuffer.KeyUp(key, modifier)
		}
	} else if action != glfw.Release {
		keyName := glfw.GetKeyName(glfwKey, scancode)
		if key, knownKey = input.ResolveShortcut(keyName, modifier); knownKey {
			window.CallKey(key, modifier)
		}
	}
}

func (window *OpenGLWindow) onChar(rawWindow *glfw.Window, char rune) {
	window.CallCharCallback(char)
}

func (window *OpenGLWindow) mapModifier(mods glfw.ModifierKey) input.Modifier {
	modifier := input.ModNone

	if (mods & glfw.ModControl) != 0 {
		modifier = modifier.With(input.ModControl)
	}
	if (mods & glfw.ModShift) != 0 {
		modifier = modifier.With(input.ModShift)
	}
	if (mods & glfw.ModAlt) != 0 {
		modifier = modifier.With(input.ModAlt)
	}
	if (mods & glfw.ModSuper) != 0 {
		modifier = modifier.With(input.ModSuper)
	}

	return modifier
}

func (window *OpenGLWindow) onDrop(rawWindow *glfw.Window, filePaths []string) {
	window.CallFileDropCallback(filePaths)
}
