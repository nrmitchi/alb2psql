class Alb2psql < Formula
  desc "Tool to fetch ALB request logs, and load into a local PostgreSQL instance"
  homepage "https://github.com/nrmitchi/alb2psql"
  url "https://github.com/nrmitchi/alb2psql/archive/v0.0.1.tar.gz"
  sha256 "2c7da89577182db466e60f8c3d6b2f35e48d2e876435f395c10e2e49facbef7f"
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
