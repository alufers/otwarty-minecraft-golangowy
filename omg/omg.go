package omg

import (
	"fmt"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const (
	windowWidth  = 960
	windowHeight = 540
)

func Main() {
	runtime.LockOSThread()

	if err := glfw.Init(); err != nil {
		panic(fmt.Errorf("could not initialize glfw: %v", err))
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	win, err := glfw.CreateWindow(800, 600, "otwarty minecraft golangowy", nil, nil)

	if err != nil {
		panic(fmt.Errorf("could not create opengl renderer: %v", err))
	}

	win.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		panic(err)
	}

	gl.ClearColor(0, 0.5, 1.0, 1.0)

	shader := newShaderProgram(`#version 410 core
	#extension GL_ARB_separate_shader_objects : enable
	#extension GL_ARB_shading_language_420pack : enable
	layout (location = 0) in vec3 aPos;
layout (location = 1) in vec3 aColor;
layout (location = 2) in vec2 aTexCoord;

out vec3 ourColor;
out vec2 TexCoord;

void main()
{
    gl_Position = vec4(aPos, 1.0);
    ourColor = aColor;
    TexCoord = aTexCoord;
}`, `#version 410 core
	out vec4 FragColor;
	in vec3 ourColor;
	in vec2 TexCoord;

	uniform sampler2D albedo;

	void main()
	{
		FragColor = texture(albedo, TexCoord);
	}`)
	err = shader.compile()

	if err != nil {
		panic(err)
	}

	m := &mesh{
		shaderProgram: shader,
		vertices: []float32{
			// positions          // colors           // texture coords
			0.5, 0.5, 0.0, 1.0, 0.0, 0.0, 1.0, 1.0, // top right
			0.5, -0.5, 0.0, 0.0, 1.0, 0.0, 1.0, 0.0, // bottom right
			-0.5, -0.5, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, // bottom left
			-0.5, 0.5, 0.0, 1.0, 1.0, 0.0, 0.0, 1.0, // top left
		},
		indices: []uint32{
			0, 1, 3, // first triangle
			1, 2, 3, // second triangle
		},
	}

	m.bindBuffers()

	t := newTexture("textures/grass_side.png")
	err = t.load()
	if err != nil {
		panic(err)
	}
	for !win.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		shader.setTexture(uniformAlbedo, t)
		m.draw()

		win.SwapBuffers()
		glfw.PollEvents()
	}
	m.dispose()
	glfw.Terminate()
}
