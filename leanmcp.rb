class Leanmcp < Formula
  desc "LeanMCP CLI - Manage projects and chats"
  homepage "https://github.com/rosaboyle/leanmcp-cli-chat-deploy"
  version "1.0.1"

  if OS.mac? && Hardware::CPU.arm?
    url "https://github.com/rosaboyle/leanmcp-cli-chat-deploy/releases/download/v1.0.1/leanmcp-1.0.1-darwin-arm64.tar.gz"
    sha256 "d429514b0c502e672d109218eede01f6cb432b636c0f1e7d8a2a31a538a24639"
  elsif OS.mac? && Hardware::CPU.intel?
    url "https://github.com/rosaboyle/leanmcp-cli-chat-deploy/releases/download/v1.0.1/leanmcp-1.0.1-darwin-amd64.tar.gz"
    sha256 "d429514b0c502e672d109218eede01f6cb432b636c0f1e7d8a2a31a538a24639"
  end

  def install
    bin.install "leanmcp-darwin-arm64" => "leanmcp" if OS.mac? && Hardware::CPU.arm?
    bin.install "leanmcp-darwin-amd64" => "leanmcp" if OS.mac? && Hardware::CPU.intel?
  end

  test do
    system "#{bin}/leanmcp", "version"
  end
end
