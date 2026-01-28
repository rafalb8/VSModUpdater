# Maintainer: Rafal Babinski <rafalb8@hotmail.com>

pkgname=VSModUpdater
pkgver=v1.3.0
pkgrel=1
pkgdesc='Vintage Story Mod Updater'
arch=('x86_64')
url='https://github.com/rafalb8/VSModUpdater'
license=('MIT')
makedepends=('go' 'git')

source=("git+$url.git#tag=${pkgver}")
sha256sums=('SKIP')

pkgver() {
  cd "$pkgname"
  git describe --long --tags | sed 's/\([^-]*-\)g/r\1/' | sed 's/-/./g'
}

build() {
  cd "$pkgname"
  make linux TAG="$pkgver"
}

package() {
  cd "$pkgname"
  install -D -m755 "result/$pkgname" "$pkgdir/usr/bin/$pkgname"
  install -D -m644 LICENSE "$pkgdir/usr/share/licenses/$pkgname/LICENSE"
}
