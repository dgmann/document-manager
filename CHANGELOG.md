# Changelog

## [1.12.6](https://github.com/dgmann/document-manager/compare/v1.12.5...v1.12.6) (2024-10-17)


### Bug Fixes

* **cli:** correctly close resources ([0a9d1b4](https://github.com/dgmann/document-manager/commit/0a9d1b45ec8ce73ce17b4c9caacea7d60638c29b))

## [1.12.5](https://github.com/dgmann/document-manager/compare/v1.12.4...v1.12.5) (2024-10-16)


### Bug Fixes

* **cli:** write uncategorized files to different file ([b363fb8](https://github.com/dgmann/document-manager/commit/b363fb857d4c058576decff769b133d9708efca1))

## [1.12.4](https://github.com/dgmann/document-manager/compare/v1.12.3...v1.12.4) (2024-10-16)


### Bug Fixes

* **cli:** nil pointer exception ([5e5b4a1](https://github.com/dgmann/document-manager/commit/5e5b4a11305d9800c6bdd1d20c1d16801a8ede0f))

## [1.12.3](https://github.com/dgmann/document-manager/compare/v1.12.2...v1.12.3) (2024-10-15)


### Bug Fixes

* **cli:** retry download on error ([648ac03](https://github.com/dgmann/document-manager/commit/648ac03c1cdcf6b721badecbc9f08c9045b92ba1))

## [1.12.2](https://github.com/dgmann/document-manager/compare/v1.12.1...v1.12.2) (2024-10-07)


### Bug Fixes

* **cli:** close body, now for real ([e68b978](https://github.com/dgmann/document-manager/commit/e68b97888d384f2aa27b4cf3855e29a583fe7aed))

## [1.12.1](https://github.com/dgmann/document-manager/compare/v1.12.0...v1.12.1) (2024-10-07)


### Bug Fixes

* **cli:** close response body ([4cf20b2](https://github.com/dgmann/document-manager/commit/4cf20b2fe38b23f1889f61fdfc21d140efae784c))

## [1.12.0](https://github.com/dgmann/document-manager/compare/v1.11.0...v1.12.0) (2024-10-07)


### Features

* **cli:** add page check ([a3ac754](https://github.com/dgmann/document-manager/commit/a3ac75427b83a8e942a0d26f37eea2b288e528f0))

## [1.11.0](https://github.com/dgmann/document-manager/compare/v1.10.1...v1.11.0) (2024-10-06)


### Features

* **api:** remove deleted pages from filesystem ([336dc1e](https://github.com/dgmann/document-manager/commit/336dc1e9869eb104b04263f52a4ea74e6ec0aed7))
* **cli:** add data check comand ([cb7ea46](https://github.com/dgmann/document-manager/commit/cb7ea46d00c16156485f7aa94b457c0bfe7a9309))
* enable compression ([8437fef](https://github.com/dgmann/document-manager/commit/8437fef7eb851981e6663516b905349148323e43))


### Bug Fixes

* **api:** improve deletion of records ([aaa7afd](https://github.com/dgmann/document-manager/commit/aaa7afd8925fdd6829619aba7cd9d4d2be5c3dac))
* **cli:** correct unkown patient id ([084d62a](https://github.com/dgmann/document-manager/commit/084d62a9bcd4f60a9668f80f6291844591a1d2b2))
* **m1-helper:** improve db query ([ea740e1](https://github.com/dgmann/document-manager/commit/ea740e14daa64d62af269784d8b6064e733888c6))

## [1.10.1](https://github.com/dgmann/document-manager/compare/v1.10.0...v1.10.1) (2024-10-04)


### Bug Fixes

* **cli:** correct nil pointer exception ([2d5e1fa](https://github.com/dgmann/document-manager/commit/2d5e1fa63558ac019f3bf1b357566cd4ad8ded34))

## [1.10.0](https://github.com/dgmann/document-manager/compare/v1.9.3...v1.10.0) (2024-10-03)


### Features

* **directory-watcher:** watch multiple folders ([ad5144d](https://github.com/dgmann/document-manager/commit/ad5144d4b169c9b08f1d3d413d707eca7894cac1))

## [1.9.3](https://github.com/dgmann/document-manager/compare/v1.9.2...v1.9.3) (2024-09-18)


### Bug Fixes

* **api:** set default port to 8080 ([9858a7b](https://github.com/dgmann/document-manager/commit/9858a7bb8dec7d2b9237d2cb2d40256dff39b9e5))
* **m1-adapter:** correct variable columns ([bd6f08d](https://github.com/dgmann/document-manager/commit/bd6f08d53e0eb22fe9381145aa9d208726681b01))

## [1.9.2](https://github.com/dgmann/document-manager/compare/v1.9.1...v1.9.2) (2024-09-18)


### Bug Fixes

* **pdf-processor:** correct nil pointer error ([e9b2412](https://github.com/dgmann/document-manager/commit/e9b241252a11c6644195ef7103bc71122e3cf917))

## [1.9.1](https://github.com/dgmann/document-manager/compare/v1.9.0...v1.9.1) (2024-09-18)


### Bug Fixes

* correct CMD ([a4d8763](https://github.com/dgmann/document-manager/commit/a4d8763ef4d8ed49acf4f82f37b98ecab633da3b))

## [1.9.0](https://github.com/dgmann/document-manager/compare/v1.8.1...v1.9.0) (2024-09-17)


### Features

* add record downloader cli ([34b2458](https://github.com/dgmann/document-manager/commit/34b24583f6cbc4e7e65977b5ad02dda1375a0ef9))
* **cli:** improve download command ([2dda1f5](https://github.com/dgmann/document-manager/commit/2dda1f5034cd1528ad9a6acc6847bec4bb6fc843))
* improve build performance ([a22d41e](https://github.com/dgmann/document-manager/commit/a22d41ebcc1db80e1f5c8e2a3c9ff8bc2d562551))


### Bug Fixes

* correct Dockerfile ([f191311](https://github.com/dgmann/document-manager/commit/f191311c6041d3ddb5b2cfee1b1066fc96b2ace6))
* correct dockerfiles ([71fd734](https://github.com/dgmann/document-manager/commit/71fd7347e5024980ed0b72aa2bee23ec96baa6fd))
* include ocr service in build ([a1b1133](https://github.com/dgmann/document-manager/commit/a1b11334965ee4e30c7cc23b4d3b84ee2e063a3f))

## [1.8.1](https://github.com/dgmann/document-manager/compare/v1.8.0...v1.8.1) (2024-05-13)


### Bug Fixes

* **deploy:** correct version in released docker-compose.yaml ([46fb0e7](https://github.com/dgmann/document-manager/commit/46fb0e730f1c2689d8be7a4a3846bbf29a96ccf7))

## [1.8.0](https://github.com/dgmann/document-manager/compare/v1.7.1...v1.8.0) (2024-05-13)


### Features

* **m1-adapter:** enable fuzzy search ([a7826b9](https://github.com/dgmann/document-manager/commit/a7826b967c99e3fc3a508d540f4db1eb9dcb4d15))


### Bug Fixes

* automatically build artifacts ([a593c58](https://github.com/dgmann/document-manager/commit/a593c58e2dff0aabf8f228a1c61a9c0d91484015))
* **m1-adapter:** change similarity to &gt;= ([ab1d3f7](https://github.com/dgmann/document-manager/commit/ab1d3f75ef214eaf7077b6b0dfbd1e091da477bb))

## [1.8.0](https://github.com/dgmann/document-manager/compare/v1.7.1...v1.8.0) (2024-05-13)


### Features

* **m1-adapter:** enable fuzzy search ([a7826b9](https://github.com/dgmann/document-manager/commit/a7826b967c99e3fc3a508d540f4db1eb9dcb4d15))


### Bug Fixes

* automatically build artifacts ([a593c58](https://github.com/dgmann/document-manager/commit/a593c58e2dff0aabf8f228a1c61a9c0d91484015))
* **m1-adapter:** change similarity to &gt;= ([ab1d3f7](https://github.com/dgmann/document-manager/commit/ab1d3f75ef214eaf7077b6b0dfbd1e091da477bb))

## [1.7.1](https://github.com/dgmann/document-manager/compare/v1.7.0...v1.7.1) (2024-05-13)


### Bug Fixes

* **m1-adapter:** filter patients without id ([5237fc5](https://github.com/dgmann/document-manager/commit/5237fc5967cbf060a140cdbe6132ea8b07de1657))

## [1.7.0](https://github.com/dgmann/document-manager/compare/v1.6.4...v1.7.0) (2024-05-13)


### Features

* add devcontainer ([4399308](https://github.com/dgmann/document-manager/commit/43993089be84125eb26dda698631331cbb6946b3))
* add File Watcher run config ([d8b21ba](https://github.com/dgmann/document-manager/commit/d8b21bafb75938a6f2fc021826158597d00140d3))
* add run config for devcontainer ([e40ed36](https://github.com/dgmann/document-manager/commit/e40ed36dc0ee3cca68e82bae2561bba59d6b1d0b))
* **devcontainer:** add app ([b7b7cfd](https://github.com/dgmann/document-manager/commit/b7b7cfd69b50a12410eaad7beaa9f176d9814c8c))
* **devcontainer:** add observability stack ([acad583](https://github.com/dgmann/document-manager/commit/acad5834f4c7789b8b0c244a275a1fee135f2016))
* **devcontainer:** use prebuilt image ([2f13091](https://github.com/dgmann/document-manager/commit/2f130914ccc3a83158d1ebca4a258b490d400c61))
* **directory-watcher:** migrate to api client package ([51c03b9](https://github.com/dgmann/document-manager/commit/51c03b976b147c5da771fe27f3bc746c123f3dc7))
* improve otel logging ([bd1fa76](https://github.com/dgmann/document-manager/commit/bd1fa76641f6e75125ae218550184bb2724b5125))
* **m1-adapter:** switch to pure go driver ([e11fc85](https://github.com/dgmann/document-manager/commit/e11fc8563d210fdcb9ed233a860f0b9e0cff8393))


### Bug Fixes

* **devcontainer:** add missing dependencies ([e55fe0f](https://github.com/dgmann/document-manager/commit/e55fe0f01fd4f6c790891c51592351ac2583cdd4))
* **frontend:** adapt run config ([9a081d7](https://github.com/dgmann/document-manager/commit/9a081d7fb8a17714790db8d230e052c3968daee9))
* **frontend:** fix tests ([b87ca84](https://github.com/dgmann/document-manager/commit/b87ca847e25c023cad6f2bf7bbdd75f85a3d2a67))
* **m1-helper:** fix tests ([48bd643](https://github.com/dgmann/document-manager/commit/48bd643360d11eb41fcb6343138c74447e25899f))
* **m1-helper:** update go.(mod|sum) ([7aa1e16](https://github.com/dgmann/document-manager/commit/7aa1e16d3d778a428c89b34e638719530c3d2056))
* **ocr:** fix tests ([d4988eb](https://github.com/dgmann/document-manager/commit/d4988eb2d4b405439e9621ef43e473cf93a71b4a))
