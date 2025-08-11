class Mcli < Formula
  desc "LeanMCP CLI - Manage projects and chats"
  homepage "https://github.com/rosaboyle/leanmcp-cli-chat-deploy"
  version "1.0.1"

  if OS.mac? && Hardware::CPU.arm?
    url "https://github.com/rosaboyle/leanmcp-cli-chat-deploy/releases/download/v1.0.1/mcli-1.0.1-darwin-arm64.tar.gz"
    sha256 "d429514b0c502e672d109218eede01f6cb432b636c0f1e7d8a2a31a538a24639"
  elsif OS.mac? && Hardware::CPU.intel?
    url "https://github.com/rosaboyle/leanmcp-cli-chat-deploy/releases/download/v1.0.1/mcli-1.0.1-darwin-arm64.tar.gz"
    sha256 "d429514b0c502e672d109218eede01f6cb432b636c0f1e7d8a2a31a538a24639"
  end

  def install
    bin.install "mcli-darwin-arm64" => "mcli" if OS.mac? && Hardware::CPU.arm?
    bin.install "mcli-darwin-amd64" => "mcli" if OS.mac? && Hardware::CPU.intel?
  end

  test do
    system "#{bin}/mcli", "version"
  end
end
