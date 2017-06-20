class Alb2psql < Formula
  desc "Tool to fetch ALB request logs, and load into a local PostgreSQL instance"
  homepage "https://github.com/nrmitchi/alb2psql"
  url "https://github.com/nrmitchi/alb2psql/archive/v0.0.2.tar.gz"
  sha256 "728eea38bbed522a12013e9c9228223fcc296867e7107ace745f5bb4362cb111"
  head "https://github.com/nrmitchi/alb2psql.git", :branch => "master"
  bottle :unneeded

  # option "with-short-names", "link as \"kctx\" and \"kns\" instead"

  def install
    bin.install "alb2psql" => "alb2psql"
    # include.install "utils.bash"

    # bash_completion.install "completion/kubectx.bash" => "kubectx"
    # bash_completion.install "completion/kubens.bash" => "kubens"
    # zsh_completion.install "completion/kubectx.zsh" => "_kubectx"
    # zsh_completion.install "completion/kubens.zsh" => "_kubens"
  end

  test do
    system "which", "alb2psql"
  end
end
