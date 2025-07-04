## unreleased
- build配下にリネーム漏れのファイルが存在したのを修正

## v0.4.15-alpha1-cs1.1.1
- Updated references from "CasaOS" to "CassetteOS" in comments and documentation
- Bumped cassetteos-common module version to include related reference updates
## v0.4.15-alpha1-cs1.1.0
### Added
- Add GitHub Action `ci.yml` to run tests on push to `main` and `develop` branches.

### Changed
- Renamed project references from **CasaOS-Gateway** to **CassetteOS-Gateway** (including module path, binary names, README, and internal scripts).
- Updated `.goreleaser.yaml` and `.goreleaser.debug.yaml` to reflect new naming and embed build metadata (`commit`, `date`, `version`) via `ldflags`.
- Revised migration and setup scripts under `build/scripts` to align with CassetteOS context.
- Customized GitHub Actions workflow (`release.yml`) for simplified GoReleaser execution under CassetteOS namespace.

## v0.4.15-alpha1-cs1.0.0
- Based on CasaOS v0.4.15
- - Replaced module paths to use our own GitHub fork instead of the original IceWhaleTech repository.
  (e.g., `github.com/IceWhaleTech/CasaOS-Gateway` → `github.com/BeesNestInc/CassetteOS-Gateway`)
