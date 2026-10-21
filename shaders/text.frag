// this shader is using in Go like embed.
#version 410 core
in vec2 TexCoord;
out vec4 FragColor;
uniform sampler2D uFontTexture; // text via font
uniform vec4 uTextColor; // text color (R,G,B)
void main() {
	vec4 texColor = texture(uFontTexture, TexCoord);
	if (texColor.a < 0.1) {
		discard;
	}
	// dot color texture to color text (saving alpha-channel)
	FragColor = vec4(uTextColor.rgb, texColor.a);
}
