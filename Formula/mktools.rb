# Formula/mktools.rb
class Mktools < Formula
  desc "CLI tool for generating LLM context and development tasks"
  homepage "https://github.com/amenophis1er/mktools"
  version "v0.1.0" # Automatically updated by release workflow

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/amenophis1er/mktools/releases/download/#{version}/mktools-darwin-arm64"
      sha256 "..." # darwin-arm64
    else
      url "https://github.com/amenophis1er/mktools/releases/download/#{version}/mktools-darwin-amd64"
      sha256 "..." # darwin-amd64
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/amenophis1er/mktools/releases/download/#{version}/mktools-linux-arm64"
      sha256 "..." # linux-arm64
    else
      url "https://github.com/amenophis1er/mktools/releases/download/#{version}/mktools-linux-amd64"
      sha256 "..." # linux-amd64
    end
  end

  def install
    binary_name = "mktools"
    bin.install Dir["mktools-*"][0] => binary_name
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/mktools version")
  end
end