# Maintainer: Kris Nóva <kris@nivenly.com>

pkgbase=naml
pkgname=(naml)
pkgver=v1.0.3
pkgrel=1
pkgdesc="Not Another Markup Language [NAML]"
arch=(x86_64)
url="https://github.com/kris-nova/naml"
license=(Apache)
makedepends=(go make)
checkdepends=()
optdepends=()
backup=()
options=()
source=("git+https://github.com/kris-nova/naml.git")
sha256sums=('SKIP')

# Kris Nóva PGP Key
#validpgpkeys=('F5F9B56417B7F2CAC1DEC2E372BB115B4DDD8252')

build() {
	cd $pkgname
	git checkout tags/$pkgver -b $pkgver
	GO111MODULE=on
	go mod vendor
	go mod download
	make
}

package() {
	cd $pkgname
	make install
}
