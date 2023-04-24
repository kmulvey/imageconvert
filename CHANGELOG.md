
2.3.0
=============
2023-03-02

* bump testify (2d158aa1)
* upgrade path lib (1cefdecd)
* normalize github actions (711ee25b)
* normalize goreleaser config (307ae5c3)
* increace dedup timeout (0381946e)
* bump deps to get latest path (908a9c77)
* start supporting resizing (1230ca67)
* lust log dir name (dd691277)
* ignore vscode (0a714eca)
* build targest and clean (7ff8f534)
* bump path (533c1af9)
* fix goreleaser warnings (031d1814)
* bump path (3eb26a44)
* remove verbose output (ad4b6e77)
* update path and rand (03dfe07e)
* bump deps (438450a3)
* deps (e6a04a3f)
* changes required for new path lib version (c960c0db)
* change how we quit so defers get exececuted (543f59e3)
* spinner for trimlog (06579db5)
* need long timeouts for large file writes (0b986fe1)
* fix filenames with weird characters (ce3b26a1)
* deps for watch (2f01f702)
* ignore binaries (b28f3642)
* better watch comment (0bb6ba38)
* break out some code from main() (ddb168e0)
* reorganize main() so its more understandable (d8d2a043)
* clean up main() func (c6c7b445)
* cleanup and comments (54c495f7)
* more binaries to build (158dfb2a)
* fix watch tests (402de72d)
* init trimlog (6b4c4387)
* beginnings of watch support (14a131c8)
* build all binaries (da77c5c2)
* unify build configs (062135e3)
* Merge pull request #16 from kmulvey/dependabot/go_modules/golang.org/x/image-0.2.0 (fceb5ee3)
* Bump golang.org/x/image from 0.1.0 to 0.2.0 (09835ea7)
* upgrade path lib to correct time ranges (ad95fa93)
* simplify logging and include progress indicator (0a3db60e)
* correct  typo (d8271f17)
* wider regex to rename files (d80baa9e)
* change modified-since flag to be more clear (a13dc09f)
* deps (e5656948)
* remove cask (c2713fa2)
* cross platform install (db4a35ab)
* install deps (d06dbd6e)
* lint input (40be74d0)
* lint (252da4e0)
* TestCompressJPEG (35eaa86e)
* makeTestDir (9bf5be5b)
* test file escape (ea3ced14)
* test identify (5b1e8032)
* move convert tests (fc5b77dc)
* random dirs for test (81f2b882)
* correct extension regex (5308cd77)
* use real path (d2bd1036)
* deps (0c2eecc8)
* true jpg test image (91244578)
* test convert errors cases (e606a5a9)
* rename all jpegs to .jpg (15f484c8)
* rename list file (1de431fb)
* fix deps after rebase (a5605473)
* add fake jpg test file (c9affcc0)
* remove old code that now lives in path lib (774dfab7)
* rearrange worker for easier testing (233c4719)
* comments and remove logging from lib code (b7270a54)
* move logging to client code (23809069)
* comments, reorder removing the old file, rename variables (6626b6ae)
* move logging into client code (e09e8e5f)
* global regex so its the same every time (7a1be10d)
* upgrade path lib (697fa7bb)
* typo (8498db3f)
* refactoring (87137b9b)
* deps (d9e458d2)
* rework how collision are handled (63017383)

2.2.0
=============
2022-11-10

* deps (7065eabe)
* if->switch (11bc7cc0)
* Merge pull request #10 from kmulvey/dependabot/go_modules/github.com/kmulvey/humantime-0.4.3 (6d5d996c)
* Bump github.com/kmulvey/humantime from 0.4.2 to 0.4.3 (ba0da62e)
* Merge pull request #11 from kmulvey/dependabot/go_modules/github.com/stretchr/testify-1.8.1 (430a08ad)
* rework how rename works (970ad4fd)
* add force option to bypass processed map (e5c8f836)
* Bump github.com/stretchr/testify from 1.8.0 to 1.8.1 (dd41a141)
* Merge pull request #9 from kmulvey/dependabot/go_modules/github.com/kmulvey/path-0.9.0 (ec5c4d6f)
* Bump github.com/kmulvey/path from 0.8.0 to 0.9.0 (abaa2b48)
* Merge pull request #8 from kmulvey/dependabot/go_modules/go.szostok.io/version-1.1.0 (c2a2fda5)
* Bump go.szostok.io/version from 1.0.0 to 1.1.0 (c0537671)
* remove unused error (12438fd9)
* print help message (fba113dc)
* deps (42bad873)
* deps (dfc4e946)
* deps (f3b1b42e)
* use latest path (976ea768)
* print version (8cb75c93)
* dont quit on errors (d5e71bfc)
* deps (7ff294ed)
* use new func (19a941aa)
* Merge pull request #2 from kmulvey/dependabot/go_modules/github.com/kmulvey/humantime-0.4.2 (65039c55)
* Bump github.com/kmulvey/humantime from 0.4.1 to 0.4.2 (6c21aecb)
* linting (49f581ce)
* comments, typo, flag.PrintDefaults (3cb9d5da)
* upgrade path (0bc7610c)
* filter by regex (15a65570)
* update path lib (0ef3817e)
* linting (e95b407e)
* use path and humantime (0e8f54cf)
* deps (0339e565)
* reset modtime (e2b9a78b)
* bump deps (f606995f)
* support glob input (eb073ec3)

2.1.0
=============
2022-08-07

* badges (0cc432e0)
* lint (76b6ad5e)
* go19 and deps (fe4541f4)

2.0.0
=============
2022-07-24

* bump deps (f144cc08)
* Merge pull request #1 from kmulvey/dependabot/go_modules/github.com/sirupsen/logrus-1.9.0 (83eed292)
* try fixing releaser (0aff9845)
* add threading (fa5e2df8)
* workers (b3b60452)
* Bump github.com/sirupsen/logrus from 1.8.1 to 1.9.0 (4f926038)
* change how the file list is built (ee7375bd)
* remove all HandleErr (26030880)
* more errors (8c526128)
* remove HandleErr (9022d130)
* remove HandleErr (6d831c10)
* remove HandleErr (8721ae33)
* me (f96d6bd2)
* vuln scan (c1d6e736)
* deps (0e34909f)
* v3 actions (870c1e60)
* copy pasta (316d532a)

1.0.2
=============
2022-06-02

* deps (2756f41c)
* deps (99e019e8)
* bump deps (8f9d3a0f)

1.0.1
=============
2022-04-19

* lint timeout (f73e28d2)

1.0.0
=============
2022-04-19

* go 18 (ce5e3722)
* deps (e760036b)
* convert integration test (0e13731e)

0.4.0
=============
2022-04-13

* 18 deps (076028fd)
* deps (010c225a)
* Update README.md (d4c78e29)
* readme (b2ae3a94)
* run notes (0149a16a)
* deps (9f355c4d)
* lint (ff0121ef)

0.3.0
=============
2022-04-01

* use full file path (b3e6a7fa)
* Update README.md (deb2dfe9)

0.2.0
=============
2022-03-31

* log time format (ad65d3d8)
* batter errors (33d4a3d9)
* bail early (d98d95e7)
* handle fake jpegs (ec11e389)
* rename jpegs (51274530)
* logging (daf81e90)
* debugging totals (21e116b0)
* ul (4e0dac97)

0.1.0
=============
2022-03-30

* dont overwrite existing jpgs (b6e05fc5)
* comments (4efd9744)
* deps (ae3d5ba6)
* dont convert jpegs (065425a8)
* Create README.md (c7dc5ebc)
* 1.8 (72031db8)
* deps (4c353f04)
* deps (53a4caad)
* update deps (8faafad3)
* try 1.8 (08fdf477)
* try go 17 (6f136061)
* go releaser (9371b711)
* lint (185fa628)
* github action (0675000c)
* simplify conversion, handle misnamed files (72f35e9f)
* args (d1d56c11)
* go 17 (711a9e11)
* RW file, and init map (6ff5d135)
* create a compress log so we dont do it twice (3a2bcfd0)
* rename jpg (e1665ec3)
* compare one case (ef954765)
* better compress error handling (5bc6be7c)
* better special char escaping (ae394ea6)
* handle compress error and preserve mod time (e74ddec5)
* copypasta (8bd92454)
* go 16 (633c1dcd)
* move (4e067261)
* use slices (c0731440)
* reorg main file (15a15e74)
* deps (b072e208)
* rename (dde9786d)
* arr (8bc4df10)
* compress wip (092513bb)
* test all jpgs (6e454200)
* reorg (c48994b0)
* Merge branch 'master' into main (979ee9d9)
* Initial commit (3a106132)
* init (842da0ea)


