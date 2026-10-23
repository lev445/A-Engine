package main
import (
	"log"
	"runtime"
	"strings"
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.4/glfw"
	mgl32 "github.com/go-gl/mathgl/mgl32"
	"a-engine/core" 
	_ "embed"
)
var projection mgl32.Mat4
var orthoProjection mgl32.Mat4
// shaders (this is base shaders. this shaders is using in elf file as embed.)
//go:embed shaders/base.vert
var vertexShaderSource string
//go:embed shaders/base.frag
var fragmentShaderSource string
//go:embed shaders/text.vert
var textVertexShaderSource string
//go:embed shaders/text.frag
var textFragmentShaderSource string
// for text:
var textProgram uint32
var textRenderer *core.TextRenderer
var font *core.Font


func Init() {
	runtime.LockOSThread()
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)
	//csources := gl.Str(source + "\x00")
	sourceStr := source + "\x00"
	free, freeCstr := gl.Strs(sourceStr)
	gl.ShaderSource(shader, 1, free, nil)
	freeCstr()
	gl.CompileShader(shader)
	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLenght int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLenght)
		logStr := strings.Repeat("\x00", int(logLenght+1))
		gl.GetShaderInfoLog(shader, logLenght, nil, gl.Str(logStr))
		return 0, fmt.Errorf("%s", logStr)
	}
	return shader, nil
}

func main() {
	Init()
	err := glfw.Init()
	if err != nil {
		log.Fatalln("\033[41mERR\033[0m :: app:core: Failed to initializing GLFW.", err)
	}
	defer glfw.Terminate()
	/*
	 * Let's set context
	 */
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.Resizable, glfw.True)
	/*
	 * Let's create a window.
	 */
	 window, err := glfw.CreateWindow(1270, 420, "A+ Engine", nil, nil)
	 if err != nil {
	 	log.Fatalln("\033[41mERR\033[0m :: app:core: Failed to create a window.", err)
	 }
	 defer glfw.Terminate()
	 window.MakeContextCurrent()
	 window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	 camera := core.NewCam(core.Vec3{X: 0.0, Y: 0.0, Z: 3.0})
	 var firstMouse = true
	 var lastX = float64(1270 / 2)
	 var lastY = float64(420 / 2)
	 window.SetCursorPosCallback(func(w *glfw.Window, xpos float64, ypos float64) {
	 	if firstMouse {
	 		lastX = xpos
	 		lastY = ypos
	 		firstMouse = false
	 	}
	 	xOffset := xpos - lastX
	 	yOffset := lastY - ypos
	 	lastX = xpos
	 	lastY = ypos
	 	camera.ProcessMouseMove(float32(xOffset), float32(yOffset), true)
	 })
	 window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	 	if key == glfw.KeyEscape && action == glfw.Press {
	 		w.SetShouldClose(true)
	 	}
	 })
	 if err := gl.Init(); err != nil {
	 	log.Fatalf("\033[41mERR\033[0m :: app:main.go: Failed to initialize OpenGL: %v\n", err)
	 }
	 lastFbW, lastFbH := 0, 0
	 applyFramebufferSize := func(width, height int) {
	 	// Под Wayland окно на старте может иметь framebuffer 0x0:
	 	// не строим вырожденные матрицы из нулей.
	 	if width <= 0 || height <= 0 {
	 		return
	 	}
	 	lastFbW, lastFbH = width, height
	 	gl.Viewport(0, 0, int32(width), int32(height))
	 	projection = mgl32.Perspective(mgl32.DegToRad(45.0), float32(width)/float32(height), 0.1, 100.0)
	 	orthoProjection = mgl32.Ortho2D(0, float32(width), float32(height), 0)
	 	if textRenderer != nil {
	 		textRenderer.SetProjection(orthoProjection)
	 	}
	 }
	 window.SetFramebufferSizeCallback(func(w *glfw.Window, width, height int) {
	 	applyFramebufferSize(width, height)
	 })
	 log.Println("\033[32mOK\033[0m :: app:main.go: OpenGL version:", gl.GoStr(gl.GetString(gl.VERSION)))
	 log.Println("\033[32mBUSY\033[0m :: app:main.go: Preparing to compiling shader...")
	 vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	 if err != nil {
	     log.Fatalf("\033[41mERR\033[0m :: Shader compile failed (Vertex): %v\n", err)
	 }
	 
	 fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	 if err != nil {
	     log.Fatalf("\033[41mERR\033[0m :: Shader compile failed (Fragment): %v\n", err)
	 }
	 program := gl.CreateProgram()
	 gl.AttachShader(program, vertexShader)
	 gl.AttachShader(program, fragmentShader)
	 gl.LinkProgram(program)
	 mvpLoc := gl.GetUniformLocation(program, gl.Str("uMVP\x00"))
	 vertices := []float32{
	 	-0.5, -0.5, 0.0,
	 	 0.5, -0.5, 0.0,
	 	 0.0,  0.5, 0.0,
	 }
	 var vao, vbo uint32
	 gl.GenVertexArrays(1, &vao)
	 gl.GenBuffers(1, &vbo)
	 gl.BindVertexArray(vao)
	 gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	 gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)
	 gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))
	 gl.EnableVertexAttribArray(0)
	 // initializing time 
	 lastFrame := 0.0 // FPS
	 applyFramebufferSize(window.GetFramebufferSize())
	 textRenderer = core.NewTextRenderer(orthoProjection)
	 textVert, err := compileShader(textVertexShaderSource, gl.VERTEX_SHADER)
	 if err != nil {
	     log.Fatalf("\033[41mERR\033[0m :: Text vertex shader: %v\n", err)
	 }
	 textFrag, err := compileShader(textFragmentShaderSource, gl.FRAGMENT_SHADER)
	 if err != nil {
	     log.Fatalf("\033[41mERR\033[0m :: Text fragment shader: %v\n", err)
	 }
	 textProgram = gl.CreateProgram()
	 gl.AttachShader(textProgram, textVert)
	 gl.AttachShader(textProgram, textFrag)
	 gl.LinkProgram(textProgram)
	 textRenderer.SetProgram(textProgram)
	 textRenderer.SetupVAO()
	 font, err = core.LoadFont("/usr/share/fonts/TTF/DejaVuSans.ttf", 48.0)
	 if err != nil {
	     log.Fatalf("\033[41mERR\033[0m :: Failed to load font: %v\n", err)
	 }
	 log.Println("\033[32mOK\033[0m :: Font loaded successfully")
	 
	 model := mgl32.Ident4()
	 model = mgl32.Ident4()
	 for !window.ShouldClose() {
	 	// initializing a dt (delta-time)
	 	currentFrame := glfw.GetTime()
	 	deltaTime := float32(currentFrame - lastFrame)
	 	lastFrame = currentFrame
		// под Wayland размер может приехать после первых кадров:
		// проверяем каждый кадр и не рисуем ничего до получения размера
		if fw, fh := window.GetFramebufferSize(); fw != lastFbW || fh != lastFbH {
			applyFramebufferSize(fw, fh)
		}
		if lastFbW == 0 {
			window.SwapBuffers()
			glfw.PollEvents()
			continue
		}
	 	// intializing actions (w/a/s/d)
	 	if window.GetKey(glfw.KeyW) == glfw.Press {
	 		camera.ProcessKeyboard("FORWARD", deltaTime)
	 	}
	 	if window.GetKey(glfw.KeyS) == glfw.Press {
	 		camera.ProcessKeyboard("BACKWARD", deltaTime)
	 	}
	 	if window.GetKey(glfw.KeyA) == glfw.Press {
	 		camera.ProcessKeyboard("LEFT", deltaTime)
	 	}
	 	if window.GetKey(glfw.KeyD) == glfw.Press {
	 		camera.ProcessKeyboard("RIGHT", deltaTime)
	 	}
	 	// clear screen and start a job
	 	gl.ClearColor(0.0, 0.15, 0.0, 1.0)
	 	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	 	gl.Enable(gl.DEPTH_TEST)
	 	gl.UseProgram(program)
	 	view := camera.GetViewMatrix()
	 	mvp := projection.Mul4(view).Mul4(model)
	 	gl.UniformMatrix4fv(mvpLoc, 1, false, &mvp[0])	
	 	gl.BindVertexArray(vao)
	 	gl.DrawArrays(gl.TRIANGLES, 0, 3)
	 	gl.Disable(gl.DEPTH_TEST)
	 	gl.Enable(gl.BLEND)
	 	// render text (simple-UI)
	 	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	 	textRenderer.RenderText(font, "Hello, world!", 50.0, 50.0, 1.0, 1.0, 1.0, 1.0)
	 	gl.Enable(gl.DEPTH_TEST)
	 	window.SwapBuffers()
	 	glfw.PollEvents()
	 }
	 // clean memory
	 gl.DeleteVertexArrays(1, &vao)
	 gl.DeleteBuffers(1, &vbo)
	 gl.DeleteProgram(program)
	 textRenderer.Cleanup()
	 font.Cleanup()
	 gl.DeleteProgram(textProgram)
}
