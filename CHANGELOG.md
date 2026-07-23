# Changelog

## [0.23.0](https://github.com/langchain-ai/langsmith-go/compare/v0.22.0...v0.23.0) (2026-07-23)


### Features

* **backend:** Allow storing grid layout to custom dashboards ([6cd7821](https://github.com/langchain-ai/langsmith-go/commit/6cd78210828a320fd22d2220006fd3518ad8a5b3))
* **backend:** update CRUD endpoint to allow retrieving / storing chart series metadata ([6d960a4](https://github.com/langchain-ai/langsmith-go/commit/6d960a4208323ce8c221177e3d307f52ec446de0))
* **host:** p50/p99 run queue wait-time chart in project monitoring [LSD-1547] ([986d99d](https://github.com/langchain-ai/langsmith-go/commit/986d99dde7f441bf7a85f1a7f8aa35137832e52e))
* **sandboxes:** free-form labels on sandboxes and snapshots ([65a3bc8](https://github.com/langchain-ai/langsmith-go/commit/65a3bc8d5df8e432eec9fce6f561c2cddbc48164))
* **sandboxes:** let proxy rules contribute sandbox env vars ([bbc4a24](https://github.com/langchain-ai/langsmith-go/commit/bbc4a24839a95049789dca37971a7ba23c9ba1ae))


### Bug Fixes

* **evaluators:** persist and wire through playground_settings_id for llm-as-judge oauth models [part of ENT-1402] ([9085eb7](https://github.com/langchain-ai/langsmith-go/commit/9085eb798ae6ec52537c909ae713ecfc3a42a777))
* **runs:** require session in POST /runs/stats ([beaef39](https://github.com/langchain-ai/langsmith-go/commit/beaef392cb387fc31e3ccf461e315c4872cd2f19))
* support last_queued_at in v2 runs query ([38917cb](https://github.com/langchain-ai/langsmith-go/commit/38917cbb47f616a5646095d22274aea8b4967b47))


### Reverts

* **runs:** require session in POST /runs/stats ([ec94123](https://github.com/langchain-ai/langsmith-go/commit/ec94123c1540af5cba9fb30ed7b069e6e94208ad))


### Refactors

* **runs:** remove SmithDB v2 endpoint flag [LSO-3391] ([94df746](https://github.com/langchain-ai/langsmith-go/commit/94df746124337d369cad60ef0e69121d6484512a))

## [0.22.0](https://github.com/langchain-ai/langsmith-go/compare/v0.21.1...v0.22.0) (2026-07-20)


### Features

* **runs:** add GET /v2/runs/{run_id}/url endpoint and SDK runs.get_url method ([1096438](https://github.com/langchain-ai/langsmith-go/commit/10964388b624719b5a74a0ad1039dec8e03c4f2a))

## [0.21.1](https://github.com/langchain-ai/langsmith-go/compare/v0.21.0...v0.21.1) (2026-07-17)


### Chores

* deprecated Emit() to String() ([#166](https://github.com/langchain-ai/langsmith-go/issues/166)) ([1a2b73a](https://github.com/langchain-ai/langsmith-go/commit/1a2b73a9a18d363475ed60c8ba46c3f3d850264c))

## [0.21.0](https://github.com/langchain-ai/langsmith-go/compare/v0.20.2...v0.21.0) (2026-07-17)


### Features

* **runs/v2:** make start_time optional in GET /v2/runs/:id ([18155e7](https://github.com/langchain-ai/langsmith-go/commit/18155e7b00d755d98584bd7191d2ae857e0f2f00))


### Documentation

* minor rewording of List Commits endpoint description ([f3c5e9c](https://github.com/langchain-ai/langsmith-go/commit/f3c5e9ce4c354c87c9b2390d5314e835a8774781))

## [0.20.2](https://github.com/langchain-ai/langsmith-go/compare/v0.20.1...v0.20.2) (2026-07-17)


### Chores

* new unshare endpoint that doesn't require lookup ([2908774](https://github.com/langchain-ai/langsmith-go/commit/29087746fc3d5577c670efe396568e697b29bd93))

## [0.20.1](https://github.com/langchain-ai/langsmith-go/compare/v0.20.0...v0.20.1) (2026-07-16)


### Bug Fixes

* parse large SSE events for response usage ([b1221bd](https://github.com/langchain-ai/langsmith-go/commit/b1221bdd31b7986959bd1c0a41a9cce79de78e09))
* parse large SSE events for response usage ([#160](https://github.com/langchain-ai/langsmith-go/issues/160)) ([2f60fd3](https://github.com/langchain-ai/langsmith-go/commit/2f60fd3acaa05c5da2bf3a33790ea5af5fc73c60))

## [0.20.0](https://github.com/langchain-ai/langsmith-go/compare/v0.19.0...v0.20.0) (2026-07-15)


### Features

* **context-hub:** fire commit webhook and deliver async ([a9b00f3](https://github.com/langchain-ai/langsmith-go/commit/a9b00f35b5127f025b59869696be3558e2925ded))
* **runs-stats:** forward reference_example and reference_dataset_id to SmithDB ([91020ca](https://github.com/langchain-ai/langsmith-go/commit/91020ca86923146eb2da329733ca8cb527b6b695))
* **sandboxes:** add preserve_memory_on_stop create attribute ([e6f6758](https://github.com/langchain-ai/langsmith-go/commit/e6f675879d8f4830e440a641762e019152b012ca))


### Chores

* update smithdb proto generated code ([ae90fd5](https://github.com/langchain-ai/langsmith-go/commit/ae90fd5e88ee075fcfd39ef63e3f960807b1a780))

## [0.19.0](https://github.com/langchain-ai/langsmith-go/compare/v0.18.2...v0.19.0) (2026-07-14)


### Features

* **dashboards:** convert prebuilts to v2 charts on clone [LSO-3306, LSO-3322] ([6632d79](https://github.com/langchain-ai/langsmith-go/commit/6632d7943e646d8830f17d8c06fcb62731129e60))
* **evaluators:** expose is_managed on Go evaluator reads [LSE-2532] ([0fca76d](https://github.com/langchain-ai/langsmith-go/commit/0fca76d59cfeb81a2c039235ca449c10981f78e8))


### Bug Fixes

* **smith-go:** default time bounds for POST /v2/threads/query ([ef3b87d](https://github.com/langchain-ai/langsmith-go/commit/ef3b87d536d8a61dc8d493e18d4f45571e3c7fc7))
* **smith-sdks:** name thread stats response model ThreadStats ([9568867](https://github.com/langchain-ai/langsmith-go/commit/9568867746c5c56966427e1f71126226951867cf))
* update golang.org/x/crypto to v0.52.0 ([#76](https://github.com/langchain-ai/langsmith-go/issues/76)) ([34adfef](https://github.com/langchain-ai/langsmith-go/commit/34adfefe2a1a220791ee4cb79581b73d6bde0058))


### Refactors

* **sandboxes:** make context hub sync a generic mount ([cb4c9ae](https://github.com/langchain-ai/langsmith-go/commit/cb4c9aec96804743c8a5470c273359163eea515c))


### Build System

* **deps:** bump golang.org/x/crypto from 0.51.0 to 0.52.0 in the go_modules group across 1 directory ([#75](https://github.com/langchain-ai/langsmith-go/issues/75)) ([4863ff3](https://github.com/langchain-ai/langsmith-go/commit/4863ff368afa0437e8183bf4e1ae024343bfa311))

## [0.18.2](https://github.com/langchain-ai/langsmith-go/compare/v0.18.1...v0.18.2) (2026-07-10)


### Bug Fixes

* explicit profile replaces current_profile (tenant, base URL, auth) ([#72](https://github.com/langchain-ai/langsmith-go/issues/72)) ([73b42f1](https://github.com/langchain-ai/langsmith-go/commit/73b42f17d22ab3a2ebfbd51493851dc6d1bc60e0))

## [0.18.1](https://github.com/langchain-ai/langsmith-go/compare/v0.18.0...v0.18.1) (2026-07-08)


### Documentation

* **runs:** document that limit does not apply when querying by trace ID ([1b466fe](https://github.com/langchain-ai/langsmith-go/commit/1b466fe2b31d7576b01528dcc63a413237f0880b))

## [0.18.0](https://github.com/langchain-ai/langsmith-go/compare/v0.17.0...v0.18.0) (2026-07-08)


### ⚠ BREAKING CHANGES

* **stainless:** skip sessions resource for python and typescript

### Features

* **managed-eval:** expose is_managed_evaluator on run-rule read/write paths ([74c4176](https://github.com/langchain-ai/langsmith-go/commit/74c41769bf0a10073edf9e9c732d92811ca6e67b))
* **smith-sdks:** expose threads and traces v2 endpoints in public SDK ([138b727](https://github.com/langchain-ai/langsmith-go/commit/138b727751f263b8e1d2d45b6b23ce077db4d18d))
* **stainless:** skip sessions resource for python and typescript ([e44bc58](https://github.com/langchain-ai/langsmith-go/commit/e44bc583f1997bf1d911e76b0956c0e756e3e1ed))


### Bug Fixes

* restore title/additionalProperties on v2 RunResponse JSON fields for Stainless ([13ecdcf](https://github.com/langchain-ai/langsmith-go/commit/13ecdcfdc152f7c58ae5349894b7383f8f5db8df))
* **smith-sdks:** remove dead retrieve_thread_preview SDK mapping ([0cf8801](https://github.com/langchain-ai/langsmith-go/commit/0cf88019bd2b978e03e3d28eb18465a3c85a77a2))
* **smith-sdks:** rename dataset runs POST methods from create to query ([2fabf94](https://github.com/langchain-ai/langsmith-go/commit/2fabf9471e847bb9549c9cf5b103432c441d773d))

## [0.17.0](https://github.com/langchain-ai/langsmith-go/compare/v0.16.1...v0.17.0) (2026-07-02)


### ⚠ BREAKING CHANGES

* disallow run_count session sort

### Features

* add GET endpoint to retrieve single workspace by id [closes ENT-508] ([0ca6712](https://github.com/langchain-ai/langsmith-go/commit/0ca6712d9583954517864e2f873f2825369d8746))


### Bug Fixes

* **integration:** update test types after run model rename ([#58](https://github.com/langchain-ai/langsmith-go/issues/58)) ([fe4118d](https://github.com/langchain-ai/langsmith-go/commit/fe4118d52c270b79767e2ed9f0a72d43dec518a4))


### Chores

* **examples:** dual-write source run lookup coordinates in Python ([faf14a6](https://github.com/langchain-ai/langsmith-go/commit/faf14a64b9e116950f808657495714f1f2f8af04))
* fix linter ([ed32aec](https://github.com/langchain-ai/langsmith-go/commit/ed32aec2bc870f421637ca1a87da8543c5bf8169))
* **stainless:** rename run models in runs resource ([9762121](https://github.com/langchain-ai/langsmith-go/commit/97621217812747a5c147bee8789c550dc737f9f7))


### Refactors

* disallow run_count session sort ([fb89237](https://github.com/langchain-ai/langsmith-go/commit/fb89237602d409599a2459cf7df4d89aca8f2ece))


### Build System

* **deps:** bump the major group with 4 updates ([#61](https://github.com/langchain-ai/langsmith-go/issues/61)) ([828c3b6](https://github.com/langchain-ai/langsmith-go/commit/828c3b6684a652f518f16024b2db59c5de10d9d2))
* **deps:** bump the minor-and-patch group with 8 updates ([#62](https://github.com/langchain-ai/langsmith-go/issues/62)) ([23b20a6](https://github.com/langchain-ai/langsmith-go/commit/23b20a641877c48abd4fec09a3b89826be239ba3))

## [0.16.1](https://github.com/langchain-ai/langsmith-go/compare/v0.16.0...v0.16.1) (2026-06-29)


### Bug Fixes

* trigger release for [#137](https://github.com/langchain-ai/langsmith-go/issues/137) (openai tier fanout) ([#140](https://github.com/langchain-ai/langsmith-go/issues/140)) ([fd0878a](https://github.com/langchain-ai/langsmith-go/commit/fd0878ab741b3aef02fcc3f145afcd05198c2e30))

## [0.16.0](https://github.com/langchain-ai/langsmith-go/compare/v0.15.0...v0.16.0) (2026-06-24)


### ⚠ BREAKING CHANGES

* **sandboxes:** default S3 mount endpoint

### Features

* **runs/v2:** expose runs v2 endpoints publicly with stainless config ([116cff9](https://github.com/langchain-ai/langsmith-go/commit/116cff975d008ad78c5f1fcfda49d89a72f5ee9e))
* **stats:** add include_details param to /runs/stats for backwards compat ([76eecdc](https://github.com/langchain-ai/langsmith-go/commit/76eecdc070b482ce2bd244b90dcb2162cf4d3a5d))


### Bug Fixes

* **sandboxes:** default S3 mount endpoint ([f7b0fed](https://github.com/langchain-ai/langsmith-go/commit/f7b0fedd799b8f0a876f1674b8b9660f531dffa5))


### Chores

* **smith-sdks:** update Java SDK stainless custom-code tracking file ([6312b44](https://github.com/langchain-ai/langsmith-go/commit/6312b44ffda8579b2d79774f21a1eddcabce7acb))


### Refactors

* rename field to match prev trace retention fields ([d536413](https://github.com/langchain-ai/langsmith-go/commit/d53641363ee9aff2884518e51cccf831e5b01f9d))

## [0.15.0](https://github.com/langchain-ai/langsmith-go/compare/v0.14.2...v0.15.0) (2026-06-18)


### ⚠ BREAKING CHANGES

* **sandboxes:** add gcs bucket mounts
* **sandboxes:** add ArtifactFS git mounts
* **sandboxes:** add ArtifactFS git mounts
* **sandboxes:** add gcs bucket mounts

### Features

* **abac:** add tag_value_ids to dataset creation endpoints ([93584f8](https://github.com/langchain-ai/langsmith-go/commit/93584f846fc3e5fe8e7a84f86ca10b5f9b9283b4))
* **abac:** add tag_value_ids to dataset creation endpoints ([8266323](https://github.com/langchain-ai/langsmith-go/commit/82663235ad1f136a2a7b64bc66ecb28440898c88))
* **abac:** add tag_value_ids to prompt creation endpoints [ENT-1176] ([93584f8](https://github.com/langchain-ai/langsmith-go/commit/93584f846fc3e5fe8e7a84f86ca10b5f9b9283b4))
* **abac:** add tag_value_ids to prompt creation endpoints [ENT-1176] ([313b627](https://github.com/langchain-ai/langsmith-go/commit/313b627419446c533d52ce7935e72c59c5a94d75))
* **abac:** add tag_value_ids to tracer session creation ([93584f8](https://github.com/langchain-ai/langsmith-go/commit/93584f846fc3e5fe8e7a84f86ca10b5f9b9283b4))
* **abac:** add tag_value_ids to tracer session creation ([313b627](https://github.com/langchain-ai/langsmith-go/commit/313b627419446c533d52ce7935e72c59c5a94d75))
* Add config for standard pagination ([93584f8](https://github.com/langchain-ai/langsmith-go/commit/93584f846fc3e5fe8e7a84f86ca10b5f9b9283b4))
* Add config for standard pagination ([313b627](https://github.com/langchain-ai/langsmith-go/commit/313b627419446c533d52ce7935e72c59c5a94d75))
* Add config for standard pagination ([6cb0b0b](https://github.com/langchain-ai/langsmith-go/commit/6cb0b0b48740328b4f5ab76bfe949a8afe485c03))
* add responses/compact tracing to openai responses ([#40](https://github.com/langchain-ai/langsmith-go/issues/40)) ([fecdc69](https://github.com/langchain-ai/langsmith-go/commit/fecdc69ad53b8ae594b589b3273b0a28e90abc57))
* add responses/compact tracing to openai responses ([#40](https://github.com/langchain-ai/langsmith-go/issues/40)) ([8f412d0](https://github.com/langchain-ai/langsmith-go/commit/8f412d0930787eccff7482c656ad28ec1aabc071))
* adding time to first token for experiment metrics ([93584f8](https://github.com/langchain-ai/langsmith-go/commit/93584f846fc3e5fe8e7a84f86ca10b5f9b9283b4))
* adding time to first token for experiment metrics ([313b627](https://github.com/langchain-ai/langsmith-go/commit/313b627419446c533d52ce7935e72c59c5a94d75))
* adding time to first token for experiment metrics ([9478cd5](https://github.com/langchain-ai/langsmith-go/commit/9478cd5dd6686a22fa564c6756f6d0b95dba66c2))
* allow feedback creation to skip trace retention upgrade ([fecdc69](https://github.com/langchain-ai/langsmith-go/commit/fecdc69ad53b8ae594b589b3273b0a28e90abc57))
* allow feedback creation to skip trace retention upgrade ([38966c2](https://github.com/langchain-ai/langsmith-go/commit/38966c2c0b2356a316863e51bdff99ba42c36c8e))
* **backend:** allow retrieving and storing V2 charts in database [LSO-2799] ([93584f8](https://github.com/langchain-ai/langsmith-go/commit/93584f846fc3e5fe8e7a84f86ca10b5f9b9283b4))
* **backend:** allow retrieving and storing V2 charts in database [LSO-2799] ([5366978](https://github.com/langchain-ai/langsmith-go/commit/5366978e1c292a79c852d342cc21d33b2c5a9d8b))
* end-to-end OAuth bearer for Playground [part of ENT-760] ([93584f8](https://github.com/langchain-ai/langsmith-go/commit/93584f846fc3e5fe8e7a84f86ca10b5f9b9283b4))
* end-to-end OAuth bearer for Playground [part of ENT-760] ([69c7404](https://github.com/langchain-ai/langsmith-go/commit/69c74042fc3d99e6572001b410cea846e5fe2ff1))
* **fleet:** frontend passes typed fields (closes AB-000) ([93584f8](https://github.com/langchain-ai/langsmith-go/commit/93584f846fc3e5fe8e7a84f86ca10b5f9b9283b4))
* **fleet:** frontend passes typed fields (closes AB-000) ([313b627](https://github.com/langchain-ai/langsmith-go/commit/313b627419446c533d52ce7935e72c59c5a94d75))
* **group-stats:** back thread stats with SmithDB ([fecdc69](https://github.com/langchain-ai/langsmith-go/commit/fecdc69ad53b8ae594b589b3273b0a28e90abc57))
* **group-stats:** back thread stats with SmithDB ([b12e7f0](https://github.com/langchain-ai/langsmith-go/commit/b12e7f0edb82173a7f4b84b2c12f3e4869b1e261))
* **hub:** add include_owners to repos list for Fleet [closes AB-2537] ([93584f8](https://github.com/langchain-ai/langsmith-go/commit/93584f846fc3e5fe8e7a84f86ca10b5f9b9283b4))
* **hub:** add include_owners to repos list for Fleet [closes AB-2537] ([313b627](https://github.com/langchain-ai/langsmith-go/commit/313b627419446c533d52ce7935e72c59c5a94d75))
* **hub:** add include_owners to repos list for Fleet [closes AB-2537] ([4e58829](https://github.com/langchain-ai/langsmith-go/commit/4e588294406947373fce4f04bf85b846a938d923))
* make online evaluator retention opt-in (backend) ([93584f8](https://github.com/langchain-ai/langsmith-go/commit/93584f846fc3e5fe8e7a84f86ca10b5f9b9283b4))
* make online evaluator retention opt-in (backend) ([6c91ed7](https://github.com/langchain-ai/langsmith-go/commit/6c91ed7f777decc039eb0bd47f14d357e13f5d66))
* **run-rules:** per-action trace-retention control for automations ([fecdc69](https://github.com/langchain-ai/langsmith-go/commit/fecdc69ad53b8ae594b589b3273b0a28e90abc57))
* **run-rules:** per-action trace-retention control for automations [LSO-2749] ([93584f8](https://github.com/langchain-ai/langsmith-go/commit/93584f846fc3e5fe8e7a84f86ca10b5f9b9283b4))
* **run-rules:** per-action trace-retention control for automations [LSO-2749] ([647609d](https://github.com/langchain-ai/langsmith-go/commit/647609d86b2b7386bb430856339cd689d1f43210))
* **sandboxes:** add ArtifactFS git mounts ([93584f8](https://github.com/langchain-ai/langsmith-go/commit/93584f846fc3e5fe8e7a84f86ca10b5f9b9283b4))
* **sandboxes:** add ArtifactFS git mounts ([2da0446](https://github.com/langchain-ai/langsmith-go/commit/2da04463cc600fc4b8d5da7be3e0ca0294f61770))
* **sandboxes:** add GCP proxy auth flow ([93584f8](https://github.com/langchain-ai/langsmith-go/commit/93584f846fc3e5fe8e7a84f86ca10b5f9b9283b4))
* **sandboxes:** add GCP proxy auth flow ([33efd4b](https://github.com/langchain-ai/langsmith-go/commit/33efd4b7008aaef2e60135321534ee33e4abf39b))
* **sandboxes:** add gcs bucket mounts ([93584f8](https://github.com/langchain-ai/langsmith-go/commit/93584f846fc3e5fe8e7a84f86ca10b5f9b9283b4))
* **sandboxes:** add gcs bucket mounts ([375d6bc](https://github.com/langchain-ai/langsmith-go/commit/375d6bc7ff0e91da3f903a7fc45c10539a6dea9b))
* **sandboxes:** add sandbox env var support ([93584f8](https://github.com/langchain-ai/langsmith-go/commit/93584f846fc3e5fe8e7a84f86ca10b5f9b9283b4))
* **sandboxes:** add sandbox env var support ([313b627](https://github.com/langchain-ai/langsmith-go/commit/313b627419446c533d52ce7935e72c59c5a94d75))
* **sandboxes:** filter lists by creator [INF-1492] ([93584f8](https://github.com/langchain-ai/langsmith-go/commit/93584f846fc3e5fe8e7a84f86ca10b5f9b9283b4))
* **sandboxes:** filter lists by creator [INF-1492] ([313b627](https://github.com/langchain-ai/langsmith-go/commit/313b627419446c533d52ce7935e72c59c5a94d75))
* **sandboxes:** set AWS proxy compatibility env ([93584f8](https://github.com/langchain-ai/langsmith-go/commit/93584f846fc3e5fe8e7a84f86ca10b5f9b9283b4))
* **sandboxes:** set AWS proxy compatibility env ([313b627](https://github.com/langchain-ai/langsmith-go/commit/313b627419446c533d52ce7935e72c59c5a94d75))
* **sandboxes:** snapshot memory from stopped boxes + honor restore_memory in v2 [INF-0000] ([93584f8](https://github.com/langchain-ai/langsmith-go/commit/93584f846fc3e5fe8e7a84f86ca10b5f9b9283b4))
* **sandboxes:** snapshot memory from stopped boxes + honor restore_memory in v2 [INF-0000] ([313b627](https://github.com/langchain-ai/langsmith-go/commit/313b627419446c533d52ce7935e72c59c5a94d75))
* **sandboxes:** snapshot memory from stopped boxes + honor restore_memory in v2 [INF-0000] ([d8a32fc](https://github.com/langchain-ai/langsmith-go/commit/d8a32fcaf5f3406237c42f956d7eab15e8e9b07f))


### Bug Fixes

* **evaluators:** show code evaluator trace counts [LSE-2359] ([93584f8](https://github.com/langchain-ai/langsmith-go/commit/93584f846fc3e5fe8e7a84f86ca10b5f9b9283b4))
* **evaluators:** show code evaluator trace counts [LSE-2359] ([b54b94c](https://github.com/langchain-ai/langsmith-go/commit/b54b94cc100374df31ffc56f7796567f98bdc36c))
* **sandboxes:** use built-in gcp proxy host matching ([15bfbfa](https://github.com/langchain-ai/langsmith-go/commit/15bfbfae965a1814da0262fdcafb6cab5c252cb8))
* **sandboxes:** use built-in gcp proxy host matching ([7f1026d](https://github.com/langchain-ai/langsmith-go/commit/7f1026d665dd0223407c46f6ff8ac0a429f76f56))


### Chores

* add Release Please GitHub Actions workflow ([#18](https://github.com/langchain-ai/langsmith-go/issues/18)) ([93584f8](https://github.com/langchain-ai/langsmith-go/commit/93584f846fc3e5fe8e7a84f86ca10b5f9b9283b4))
* add Release Please GitHub Actions workflow ([#18](https://github.com/langchain-ai/langsmith-go/issues/18)) ([313b627](https://github.com/langchain-ai/langsmith-go/commit/313b627419446c533d52ce7935e72c59c5a94d75))
* add Release Please GitHub Actions workflow ([#18](https://github.com/langchain-ai/langsmith-go/issues/18)) ([b5d5135](https://github.com/langchain-ai/langsmith-go/commit/b5d513516e727592fe9a6fd948b9c3a3e6250809))
* **fleet:** add param to list threads [closes AB-2522] ([93584f8](https://github.com/langchain-ai/langsmith-go/commit/93584f846fc3e5fe8e7a84f86ca10b5f9b9283b4))
* **fleet:** add param to list threads [closes AB-2522] ([313b627](https://github.com/langchain-ai/langsmith-go/commit/313b627419446c533d52ce7935e72c59c5a94d75))
* revert "Build SDK" ([#24](https://github.com/langchain-ai/langsmith-go/issues/24)) ([93584f8](https://github.com/langchain-ai/langsmith-go/commit/93584f846fc3e5fe8e7a84f86ca10b5f9b9283b4))
* revert "Build SDK" ([#24](https://github.com/langchain-ai/langsmith-go/issues/24)) ([313b627](https://github.com/langchain-ai/langsmith-go/commit/313b627419446c533d52ce7935e72c59c5a94d75))
* revert "Build SDK" ([#24](https://github.com/langchain-ai/langsmith-go/issues/24)) ([6027dd5](https://github.com/langchain-ai/langsmith-go/commit/6027dd5dd60c5d16cdbb92b2a4fde39f0e2d7c29))
* SDK release (next → main) ([#4](https://github.com/langchain-ai/langsmith-go/issues/4)) ([93584f8](https://github.com/langchain-ai/langsmith-go/commit/93584f846fc3e5fe8e7a84f86ca10b5f9b9283b4))
* SDK release (next → main) ([#4](https://github.com/langchain-ai/langsmith-go/issues/4)) ([313b627](https://github.com/langchain-ai/langsmith-go/commit/313b627419446c533d52ce7935e72c59c5a94d75))
* standardize stlc workflow filenames and names ([#37](https://github.com/langchain-ai/langsmith-go/issues/37)) ([93584f8](https://github.com/langchain-ai/langsmith-go/commit/93584f846fc3e5fe8e7a84f86ca10b5f9b9283b4))
* standardize stlc workflow filenames and names ([#37](https://github.com/langchain-ai/langsmith-go/issues/37)) ([23d09e6](https://github.com/langchain-ai/langsmith-go/commit/23d09e6b2f442e4e6fdeac7bddb47b655cf564b3))


### Refactors

* **playground:** remove legacy experiment endpoints [LSO-2230] ([93584f8](https://github.com/langchain-ai/langsmith-go/commit/93584f846fc3e5fe8e7a84f86ca10b5f9b9283b4))
* **playground:** remove legacy experiment endpoints [LSO-2230] ([313b627](https://github.com/langchain-ai/langsmith-go/commit/313b627419446c533d52ce7935e72c59c5a94d75))

## 0.14.2 (2026-06-09)

Full Changelog: [v0.14.1...v0.14.2](https://github.com/langchain-ai/langsmith-go/compare/v0.14.1...v0.14.2)

### Features

* **api:** api update ([7b3f152](https://github.com/langchain-ai/langsmith-go/commit/7b3f1526abeb72b6234bc83490e4b0c4e2e19aa7))
* **api:** api update ([9272f4a](https://github.com/langchain-ai/langsmith-go/commit/9272f4a081802422af18303c4ac399ea8e6f3404))
* **api:** api update ([80b32c4](https://github.com/langchain-ai/langsmith-go/commit/80b32c4f0b66e08a0f1896e942a1ec193f620e87))
* **api:** api update ([a20f9ad](https://github.com/langchain-ai/langsmith-go/commit/a20f9ada80cc6d1a39daee97440033e90d583172))

## 0.14.1 (2026-06-04)

Full Changelog: [v0.14.0...v0.14.1](https://github.com/langchain-ai/langsmith-go/compare/v0.14.0...v0.14.1)

### Features

* **api:** api update ([9458e67](https://github.com/langchain-ai/langsmith-go/commit/9458e67e6c3ab868ae085236ee87d026cfaa3cf0))
* **api:** api update ([d4408a1](https://github.com/langchain-ai/langsmith-go/commit/d4408a121d69a573b048c6a4e828c9910ce27a2b))
* **api:** api update ([6104c3b](https://github.com/langchain-ai/langsmith-go/commit/6104c3bc724319ad1ba83f7e51b45dfc726b1948))
* **api:** api update ([83734e4](https://github.com/langchain-ai/langsmith-go/commit/83734e47d233ab856119e9bac280d0992c11ff1a))
* **api:** api update ([aa2c032](https://github.com/langchain-ai/langsmith-go/commit/aa2c032ca6357d2a7df2f74c82a8c73b514c1143))

## 0.14.0 (2026-05-27)

Full Changelog: [v0.13.1...v0.14.0](https://github.com/langchain-ai/langsmith-go/compare/v0.13.1...v0.14.0)

### Features

* **api:** api update ([010ee31](https://github.com/langchain-ai/langsmith-go/commit/010ee3126cbefd7fff3523d065d62e243ed44977))
* **api:** api update ([995bf63](https://github.com/langchain-ai/langsmith-go/commit/995bf638bf9f8f2dc534cff2e2fe606719e49e99))

## 0.13.1 (2026-05-22)

Full Changelog: [v0.13.0...v0.13.1](https://github.com/langchain-ai/langsmith-go/compare/v0.13.0...v0.13.1)

### Bug Fixes

* openai: support `response.incomplete` status ([#99](https://github.com/langchain-ai/langsmith-go/issues/99)) ([cf3bf60](https://github.com/langchain-ai/langsmith-go/commit/cf3bf60dad4aa2fa4d8bf871b155e713bff07d84))

## 0.13.0 (2026-05-22)

Full Changelog: [v0.12.0...v0.13.0](https://github.com/langchain-ai/langsmith-go/compare/v0.12.0...v0.13.0)

### Features

* **api:** api update ([83864b4](https://github.com/langchain-ai/langsmith-go/commit/83864b4fa04c9f484347ddc0935c2178a0546ae6))


### Bug Fixes

* **auth:** add OAuth user ID header from JWT ([#98](https://github.com/langchain-ai/langsmith-go/issues/98)) ([49e4fd1](https://github.com/langchain-ai/langsmith-go/commit/49e4fd164b039c63bf5b090f74651e10c37a69be))
* **auth:** lock OAuth profile refreshes ([#97](https://github.com/langchain-ai/langsmith-go/issues/97)) ([5e6de52](https://github.com/langchain-ai/langsmith-go/commit/5e6de52a5576ea2e42ea1823f6851ffcd72050d5))

## 0.12.0 (2026-05-19)

Full Changelog: [v0.11.0...v0.12.0](https://github.com/langchain-ai/langsmith-go/compare/v0.11.0...v0.12.0)

### Features

* **api:** api update ([0c7af31](https://github.com/langchain-ai/langsmith-go/commit/0c7af31cb3b3d1ba3675e9787094d773fb901a41))
* **api:** api update ([b1fae6c](https://github.com/langchain-ai/langsmith-go/commit/b1fae6c4a9e7e3ce88c582afc78eaeb930add51e))
* **api:** api update ([47164c4](https://github.com/langchain-ai/langsmith-go/commit/47164c4ca853ed51137675839762a3c178c961aa))
* **sandbox:** default create to server runtime ([#92](https://github.com/langchain-ai/langsmith-go/issues/92)) ([89a0fae](https://github.com/langchain-ai/langsmith-go/commit/89a0fae142cd7bd51ae4571ba90b14db0ac4f23e))

## 0.11.0 (2026-05-12)

Full Changelog: [v0.10.0...v0.11.0](https://github.com/langchain-ai/langsmith-go/compare/v0.10.0...v0.11.0)

### Features

* **api:** api update ([73eda1e](https://github.com/langchain-ai/langsmith-go/commit/73eda1eba07a8472bc64f33cfd3c9cc5e90d1d2c))
* **api:** api update ([352e0fe](https://github.com/langchain-ai/langsmith-go/commit/352e0fe3320e843810574ba40961aa5a899bfdf6))
* **api:** api update ([7f92f2d](https://github.com/langchain-ai/langsmith-go/commit/7f92f2d45907613d2893850d7fb2b71d27726063))
* openai: add tool calls to `/responses` ([#88](https://github.com/langchain-ai/langsmith-go/issues/88)) ([a4a374a](https://github.com/langchain-ai/langsmith-go/commit/a4a374a798a8b5734199317fa6fa4187f81bb8b0))


### Bug Fixes

* type RunStats cost fields as float64, not string ([#90](https://github.com/langchain-ai/langsmith-go/issues/90)) ([8c3e934](https://github.com/langchain-ai/langsmith-go/commit/8c3e934887b918073ba9bd0af7c2b74394e9dc44))

## 0.10.0 (2026-05-10)

Full Changelog: [v0.9.4...v0.10.0](https://github.com/langchain-ai/langsmith-go/compare/v0.9.4...v0.10.0)

### Features

* **api:** api update ([f5cd01b](https://github.com/langchain-ai/langsmith-go/commit/f5cd01b5a28a97159266db12855770085022a6bc))
* **api:** api update ([d9ea784](https://github.com/langchain-ai/langsmith-go/commit/d9ea7844591ac1c55fff39623854a91c506276ed))
* **api:** api update ([c3a03c0](https://github.com/langchain-ai/langsmith-go/commit/c3a03c04f4b17628232169b19013d79d302fded6))
* **api:** api update ([48d267b](https://github.com/langchain-ai/langsmith-go/commit/48d267b2690a0b24ae6fde9cd22d42665d2a1b83))


### Bug Fixes

* **go:** avoid panic when http.DefaultTransport is wrapped ([f481a00](https://github.com/langchain-ai/langsmith-go/commit/f481a005dabc8fb7115a289fa805b80b80deb832))


### Chores

* redact api-key headers in debug logs ([e3add2c](https://github.com/langchain-ai/langsmith-go/commit/e3add2cfc0982c78935bb2b91e29711f0aadba41))

## 0.9.4 (2026-05-06)

Full Changelog: [v0.9.3...v0.9.4](https://github.com/langchain-ai/langsmith-go/compare/v0.9.3...v0.9.4)

### Bug Fixes

* gemini intermediate tool use ([#77](https://github.com/langchain-ai/langsmith-go/issues/77)) ([bd778e3](https://github.com/langchain-ai/langsmith-go/commit/bd778e312b7015a1254c3a9e9cb08e60838cb8b0))

## 0.9.3 (2026-05-06)

Full Changelog: [v0.10.0...v0.9.3](https://github.com/langchain-ai/langsmith-go/compare/v0.10.0...v0.9.3)

### Features

* add insights fetching to go sdk ([#40](https://github.com/langchain-ai/langsmith-go/issues/40)) ([21bbebd](https://github.com/langchain-ai/langsmith-go/commit/21bbebdf871402fb3faf11de658cacc5903dc0d1))
* added additional unit, integration, and build tests ([#35](https://github.com/langchain-ai/langsmith-go/issues/35)) ([04d6459](https://github.com/langchain-ai/langsmith-go/commit/04d6459c5b3360ba5b480b47942ca0afcc1bd305))
* added high level tracer initialization ([#21](https://github.com/langchain-ai/langsmith-go/issues/21)) ([031e9e8](https://github.com/langchain-ai/langsmith-go/commit/031e9e89a644ca7e5ca518e3c8bbec50461e0585))
* added tracing with go openai client ([#16](https://github.com/langchain-ai/langsmith-go/issues/16)) ([49d5dce](https://github.com/langchain-ai/langsmith-go/commit/49d5dce40aced07d8fef90dd4bcc5eb6391221fc))
* **api:** add go target ([07c2406](https://github.com/langchain-ai/langsmith-go/commit/07c240666a7c39e1d684090b20286dd3e5dae6ec))
* emit new_token span event on first content delta in streaming wrappers ([#67](https://github.com/langchain-ai/langsmith-go/issues/67)) ([54da463](https://github.com/langchain-ai/langsmith-go/commit/54da46347b7fb28e2c6df6aa535ba6a98790f533))
* load LangSmith profile auth in client ([#53](https://github.com/langchain-ai/langsmith-go/issues/53)) ([f15e14e](https://github.com/langchain-ai/langsmith-go/commit/f15e14e16b5e67a688b900e59ffc3e9f134c0a98))
* multipart ingestion ([#36](https://github.com/langchain-ai/langsmith-go/issues/36)) ([ff6f440](https://github.com/langchain-ai/langsmith-go/commit/ff6f440535e77ebcf805fd5a2fb589df23b135fb))
* **sandbox:** expose console command controls ([#85](https://github.com/langchain-ai/langsmith-go/issues/85)) ([d04e4ed](https://github.com/langchain-ai/langsmith-go/commit/d04e4ed15038ec6b0344f130f64a11d724167eb2))


### Bug Fixes

* added fix for missing thread token counts ([#19](https://github.com/langchain-ai/langsmith-go/issues/19)) ([6c54b0f](https://github.com/langchain-ai/langsmith-go/commit/6c54b0fbac81a2df9b0de9d38f3e3fd54180637d))


### Reverts

* load LangSmith profile auth in client ([#65](https://github.com/langchain-ai/langsmith-go/issues/65)) ([03d2f40](https://github.com/langchain-ai/langsmith-go/commit/03d2f400a3bedca1e131a6197f1e19dfed58eee3))


### Chores

* dependabot ([#26](https://github.com/langchain-ai/langsmith-go/issues/26)) ([02b58a3](https://github.com/langchain-ai/langsmith-go/commit/02b58a3fc5f16b551a0ec87b7017b0eb5f282a66))
* **deps:** bump go.opentelemetry.io/otel/sdk ([#33](https://github.com/langchain-ai/langsmith-go/issues/33)) ([675dfba](https://github.com/langchain-ai/langsmith-go/commit/675dfba4834e9de34b55580f9cb492061acb7698))
* **deps:** bump google.golang.org/grpc ([#37](https://github.com/langchain-ai/langsmith-go/issues/37)) ([bd75dd5](https://github.com/langchain-ai/langsmith-go/commit/bd75dd510aed9e68fbb4cc1ee0d7ca0673a847a1))
* rename langchain_base_url to langsmith_endpoint ([#22](https://github.com/langchain-ai/langsmith-go/issues/22)) ([66b5cb4](https://github.com/langchain-ai/langsmith-go/commit/66b5cb449050d5fff16ed86bad49fe82a99b3eee))
* stable release ([#28](https://github.com/langchain-ai/langsmith-go/issues/28)) ([127e4a7](https://github.com/langchain-ai/langsmith-go/commit/127e4a7027ff8066289a4f0407a07730209e2ee0))
* stable release ([#31](https://github.com/langchain-ai/langsmith-go/issues/31)) ([3e69fc6](https://github.com/langchain-ai/langsmith-go/commit/3e69fc60056c409c53389d5ccb2746a93eaf0107))

## 0.10.0 (2026-05-06)

Full Changelog: [v0.9.1...v0.10.0](https://github.com/langchain-ai/langsmith-go/compare/v0.9.1...v0.10.0)

### Features

* add explicit profile client option ([#82](https://github.com/langchain-ai/langsmith-go/issues/82)) ([75abe5d](https://github.com/langchain-ai/langsmith-go/commit/75abe5d418912d74058a7df7be422914a1d9bdb6))

## 0.9.1 (2026-05-06)

Full Changelog: [v0.9.0...v0.9.1](https://github.com/langchain-ai/langsmith-go/compare/v0.9.0...v0.9.1)

### Features

* **api:** api update ([2491248](https://github.com/langchain-ai/langsmith-go/commit/24912489c843b0fd600c9abab0e3a31d8e0fc561))
* **sandbox:** add Go command runtime helpers ([#79](https://github.com/langchain-ai/langsmith-go/issues/79)) ([e2d8350](https://github.com/langchain-ai/langsmith-go/commit/e2d8350b69cfde81d2f378792f0fcec81a6b5c37))

## 0.9.0 (2026-05-05)

Full Changelog: [v0.8.1...v0.9.0](https://github.com/langchain-ai/langsmith-go/compare/v0.8.1...v0.9.0)

### Features

* **sdk:** add hub directories endpoints to Stainless config ([31b6912](https://github.com/langchain-ai/langsmith-go/commit/31b6912e72e6639278333f651f4b126c6d5b75ba))

## 0.8.1 (2026-05-05)

Full Changelog: [v0.8.0...v0.8.1](https://github.com/langchain-ai/langsmith-go/compare/v0.8.0...v0.8.1)

### Features

* **api:** api update ([4797308](https://github.com/langchain-ai/langsmith-go/commit/4797308a9e8b3d61d88778e62fc99a5727ad36a4))
* **api:** api update ([83870c6](https://github.com/langchain-ai/langsmith-go/commit/83870c621fe79207ae365ad25b591a024a83dbf4))

## 0.8.0 (2026-05-04)

Full Changelog: [v0.7.0...v0.8.0](https://github.com/langchain-ai/langsmith-go/compare/v0.7.0...v0.8.0)

### Features

* **api:** api update ([8bbec4d](https://github.com/langchain-ai/langsmith-go/commit/8bbec4d9272a9435eb7ce2b6413ea6cb762c20f8))
* emit new_token span event on first content delta in streaming wrappers ([#67](https://github.com/langchain-ai/langsmith-go/issues/67)) ([54da463](https://github.com/langchain-ai/langsmith-go/commit/54da46347b7fb28e2c6df6aa535ba6a98790f533))


### Chores

* avoid embedding reflect.Type for dead code elimination ([f6f1dab](https://github.com/langchain-ai/langsmith-go/commit/f6f1dab8bdba8fc6caeef6c4ee2b834e0c7a89d6))

## 0.7.0 (2026-04-30)

Full Changelog: [v0.6.0...v0.7.0](https://github.com/langchain-ai/langsmith-go/compare/v0.6.0...v0.7.0)

### Features

* **api:** add workspaces ([ed50eeb](https://github.com/langchain-ai/langsmith-go/commit/ed50eebd3d9ac08fc1fec472bfb4a5d40a1e590c))
* **api:** api update ([b5910ea](https://github.com/langchain-ai/langsmith-go/commit/b5910eaefeba8720ff240d9178392d5d6773d307))
* **api:** manual updates ([5e08957](https://github.com/langchain-ai/langsmith-go/commit/5e089576d81cde6a023e8700b0b93e86f62cd07c))
* **api:** manual updates ([a1398c9](https://github.com/langchain-ai/langsmith-go/commit/a1398c9a11d955882bf9a10ff8448ab7bb3d163a))

## 0.6.0 (2026-04-29)

Full Changelog: [v0.5.0...v0.6.0](https://github.com/langchain-ai/langsmith-go/compare/v0.5.0...v0.6.0)

### Features

* **api:** api update ([f478bad](https://github.com/langchain-ai/langsmith-go/commit/f478bad4cb70beadb709b5318b9984cd86dd490c))
* **api:** api update ([fa2dbdb](https://github.com/langchain-ai/langsmith-go/commit/fa2dbdb60e9e570876d006655aae7556604ef05b))
* **go:** add default http client with timeout ([14b7982](https://github.com/langchain-ai/langsmith-go/commit/14b7982638dca6d1e8c057176f9a8717e03d25f8))
* load LangSmith profile auth in client ([#53](https://github.com/langchain-ai/langsmith-go/issues/53)) ([f15e14e](https://github.com/langchain-ai/langsmith-go/commit/f15e14e16b5e67a688b900e59ffc3e9f134c0a98))
* support setting headers via env ([cdb1356](https://github.com/langchain-ai/langsmith-go/commit/cdb13565cbe6016081c90e9f1cef88c89036d92e))


### Reverts

* load LangSmith profile auth in client ([#65](https://github.com/langchain-ai/langsmith-go/issues/65)) ([03d2f40](https://github.com/langchain-ai/langsmith-go/commit/03d2f400a3bedca1e131a6197f1e19dfed58eee3))


### Chores

* sync next with main for release ([#64](https://github.com/langchain-ai/langsmith-go/issues/64)) ([1458751](https://github.com/langchain-ai/langsmith-go/commit/14587513b78b4f4c81163b40c1235d8e4662d4bc))

## 0.5.0 (2026-04-23)

Full Changelog: [v0.4.0...v0.5.0](https://github.com/langchain-ai/langsmith-go/compare/v0.4.0...v0.5.0)

### Features

* **api:** api update ([7c5a612](https://github.com/langchain-ai/langsmith-go/commit/7c5a612276aa931b0a238d02f329af40c99fb19e))
* **evaluators:** add list evaluators (GET /api/v1/runs/rules) ([8ddcd54](https://github.com/langchain-ai/langsmith-go/commit/8ddcd545dca0d643137b3f7fc2bd84bf29f9bebe))


### Chores

* **internal:** more robust bootstrap script ([592b556](https://github.com/langchain-ai/langsmith-go/commit/592b5568c42e8b3de6db283e68d64676eac197f3))

## 0.4.0 (2026-04-21)

Full Changelog: [v0.3.0...v0.4.0](https://github.com/langchain-ai/langsmith-go/compare/v0.3.0...v0.4.0)

### Features

* **api:** api update ([0f99b31](https://github.com/langchain-ai/langsmith-go/commit/0f99b3146284bbd2d2963aa5d42510b2c7ce4a43))
* **api:** api update ([cd16ebc](https://github.com/langchain-ai/langsmith-go/commit/cd16ebc19568ca9f3acec14a83cdbf2c6d2a6283))
* **api:** api update ([ce8d1de](https://github.com/langchain-ai/langsmith-go/commit/ce8d1de2b5a1b8854c79edd2aa4e17df382fbe3d))
* **api:** api update ([c8717cc](https://github.com/langchain-ai/langsmith-go/commit/c8717cc4fc0868e1dd0a579ba3139fc6ec3ccbbd))
* **api:** sandbox apis ([4ef1a1e](https://github.com/langchain-ai/langsmith-go/commit/4ef1a1e741fc129d81938e24671a0db45875ca03))


### Chores

* **deps:** bump go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp ([#54](https://github.com/langchain-ai/langsmith-go/issues/54)) ([1d68865](https://github.com/langchain-ai/langsmith-go/commit/1d68865db174a5d058d70ba5b117c7267201eacd))


### Documentation

* improve examples ([38a244e](https://github.com/langchain-ai/langsmith-go/commit/38a244efc1139cdc4fc0909065acb5dd033cf240))

## 0.3.0 (2026-04-08)

Full Changelog: [v0.2.2...v0.3.0](https://github.com/langchain-ai/langsmith-go/compare/v0.2.2...v0.3.0)

### Features

* **api:** api update ([82c4094](https://github.com/langchain-ai/langsmith-go/commit/82c40941f0a12d8045d921bddfb9d05d54890cf4))
* **api:** api update ([453c395](https://github.com/langchain-ai/langsmith-go/commit/453c395d03e396da59ae9803efa6bff04693e470))
* **api:** api update ([8f3f28a](https://github.com/langchain-ai/langsmith-go/commit/8f3f28aedd17122d51cefb60a4643aa400b68872))
* **api:** api update ([02c60c1](https://github.com/langchain-ai/langsmith-go/commit/02c60c10f9c0b5fff61da3e9bd35711fa7767ee9))
* **api:** api update ([010d179](https://github.com/langchain-ai/langsmith-go/commit/010d179061be24f8abde651a84a5e747020a614a))
* **api:** api update ([c21ec37](https://github.com/langchain-ai/langsmith-go/commit/c21ec374d36a7438a5d0deeaf209464d6ef6f1b7))
* **api:** api update ([27b413d](https://github.com/langchain-ai/langsmith-go/commit/27b413de09e71dd13e5082c2e9d0532cbc4d56ba))

## 0.2.2 (2026-03-30)

Full Changelog: [v0.2.1...v0.2.2](https://github.com/langchain-ai/langsmith-go/compare/v0.2.1...v0.2.2)

### Features

* **api:** api update ([4838570](https://github.com/langchain-ai/langsmith-go/commit/483857019c4b3c4111d93b1b5f0044593b25fb82))
* **api:** api update ([bc2f46b](https://github.com/langchain-ai/langsmith-go/commit/bc2f46b211a935791f8502e0b6eaec6efb45093e))
* **api:** api update ([31b4977](https://github.com/langchain-ai/langsmith-go/commit/31b497729cfb81551e7dc903c8dd6fb6dacc61e5))
* **api:** api update ([59aa8e1](https://github.com/langchain-ai/langsmith-go/commit/59aa8e1df6df83d67779335e2881794571472707))
* **api:** manual updates ([9d68d4c](https://github.com/langchain-ai/langsmith-go/commit/9d68d4c762851ef9b42944e6a1dd5a9a5bb902cf))
* **internal:** support comma format in multipart form encoding ([8c561c7](https://github.com/langchain-ai/langsmith-go/commit/8c561c72b74ead65cbe4de6005ad8fd3424592e8))
* multipart ingestion ([#36](https://github.com/langchain-ai/langsmith-go/issues/36)) ([ff6f440](https://github.com/langchain-ai/langsmith-go/commit/ff6f440535e77ebcf805fd5a2fb589df23b135fb))


### Bug Fixes

* prevent duplicate ? in query params ([87a3062](https://github.com/langchain-ai/langsmith-go/commit/87a30622638f311349a19b4fcfa4d1f28476ae63))


### Chores

* add target-branch next for Stainless posture compliance ([#48](https://github.com/langchain-ai/langsmith-go/issues/48)) ([e190929](https://github.com/langchain-ai/langsmith-go/commit/e190929227183902197ae1a20bd267071bd9b26e))
* **ci:** skip lint on metadata-only changes ([9a00886](https://github.com/langchain-ai/langsmith-go/commit/9a00886ca5a171d0e5e8dc43e78672b9643c5a70))
* **ci:** support opting out of skipping builds on metadata-only commits ([d61b04d](https://github.com/langchain-ai/langsmith-go/commit/d61b04dfeb7be85aad1c4acbc09fc2b3250191a0))
* dependabot ([#26](https://github.com/langchain-ai/langsmith-go/issues/26)) ([02b58a3](https://github.com/langchain-ai/langsmith-go/commit/02b58a3fc5f16b551a0ec87b7017b0eb5f282a66))
* **internal:** update gitignore ([701ad97](https://github.com/langchain-ai/langsmith-go/commit/701ad973197958f06c57f0ca4a99f13d357b6a96))
* remove unnecessary error check for url parsing ([3fb1e50](https://github.com/langchain-ai/langsmith-go/commit/3fb1e50adbbda3f8e5bfbdb7ba1da9f198dd77db))

## 0.2.1 (2026-03-20)

Full Changelog: [v0.2.0...v0.2.1](https://github.com/langchain-ai/langsmith-go/compare/v0.2.0...v0.2.1)

### Features

* add insights fetching to go sdk ([#40](https://github.com/langchain-ai/langsmith-go/issues/40)) ([21bbebd](https://github.com/langchain-ai/langsmith-go/commit/21bbebdf871402fb3faf11de658cacc5903dc0d1))
* **api:** manual updates ([174f299](https://github.com/langchain-ai/langsmith-go/commit/174f299ca32d15cda0fe1a1e098afdf7349e1da6))

## 0.2.0 (2026-03-20)

Full Changelog: [v0.1.0...v0.2.0](https://github.com/langchain-ai/langsmith-go/compare/v0.1.0...v0.2.0)

### Features

* added additional unit, integration, and build tests ([#35](https://github.com/langchain-ai/langsmith-go/issues/35)) ([04d6459](https://github.com/langchain-ai/langsmith-go/commit/04d6459c5b3360ba5b480b47942ca0afcc1bd305))
* **api:** api update ([8c6ea84](https://github.com/langchain-ai/langsmith-go/commit/8c6ea84c59f26519de8df7c9f0a49f5010d2ae97))
* **api:** api update ([44b6787](https://github.com/langchain-ai/langsmith-go/commit/44b67875d837ee81272190fa2f6cc3d7d293838d))
* **api:** api update ([f9a6d01](https://github.com/langchain-ai/langsmith-go/commit/f9a6d01eb7117d07b2c389bc5498274c78a8b09c))
* **api:** api update ([21524b5](https://github.com/langchain-ai/langsmith-go/commit/21524b5ba8b7c94eb532de1f49aa285e4fd983de))
* **api:** api update ([b67882f](https://github.com/langchain-ai/langsmith-go/commit/b67882f256ae2a50924554361403a388ba134cc9))
* **api:** api update ([8f5a97b](https://github.com/langchain-ai/langsmith-go/commit/8f5a97b00167f8147040959411fb3e447c7c4b21))
* **api:** api update ([cb7bca2](https://github.com/langchain-ai/langsmith-go/commit/cb7bca2e6e8676b3915e628f57bf05568a1b3e1b))
* **api:** api update ([27c213a](https://github.com/langchain-ai/langsmith-go/commit/27c213a5bb7dac3f4f79a97ffbf6eab70d416074))
* **api:** api update ([76645a9](https://github.com/langchain-ai/langsmith-go/commit/76645a9523c50eb63caf95b1f82a4de9652dac46))
* **api:** api update ([e5b911d](https://github.com/langchain-ai/langsmith-go/commit/e5b911d68676a8d6a45137cb063c05413a225123))
* **api:** api update ([1b97b1f](https://github.com/langchain-ai/langsmith-go/commit/1b97b1f0e2e1b0153eb55c1534a7d722c58b69e7))
* **api:** api update ([84f0890](https://github.com/langchain-ai/langsmith-go/commit/84f0890c7f9ef2f85a1a2f9b609cb8c885132670))
* **api:** api update ([5d441d4](https://github.com/langchain-ai/langsmith-go/commit/5d441d4da2ae3f7d9a69af8104168525606dab75))
* **api:** manual updates ([2c1f662](https://github.com/langchain-ai/langsmith-go/commit/2c1f66280102acb26c02f004ee63ba560bbfe0a5))
* **api:** manual updates ([fc9015e](https://github.com/langchain-ai/langsmith-go/commit/fc9015e9d5267a0a1aa88ce6a313d4125df9b3f0))
* **api:** manual updates ([d8c5c5f](https://github.com/langchain-ai/langsmith-go/commit/d8c5c5f475e879a16eb275b2ba89a1e23a13696f))
* **api:** manual updates ([ee37b62](https://github.com/langchain-ai/langsmith-go/commit/ee37b62b1fa6f54cd0a0f33cda0e70cbe445788c))


### Bug Fixes

* allow canceling a request while it is waiting to retry ([1e74640](https://github.com/langchain-ai/langsmith-go/commit/1e74640a8bf7c5ed31666c737ac2e164c0e8e9d2))
* **internal:** skip tests that depend on mock server ([1a5872b](https://github.com/langchain-ai/langsmith-go/commit/1a5872b68ddb36de94355d0901b1f7ede5e76d56))
* type mismatch ([#39](https://github.com/langchain-ai/langsmith-go/issues/39)) ([bd78bb5](https://github.com/langchain-ai/langsmith-go/commit/bd78bb585c5538b289f00430921489071ce91e08))


### Chores

* change user client ([#38](https://github.com/langchain-ai/langsmith-go/issues/38)) ([76850c7](https://github.com/langchain-ai/langsmith-go/commit/76850c7cfc34065b4a42141e97b794e12e437e6e))
* **ci:** skip uploading artifacts on stainless-internal branches ([061ca81](https://github.com/langchain-ai/langsmith-go/commit/061ca81134fd8a077295a2764e9c902e055d223d))
* **deps:** bump go.opentelemetry.io/otel/sdk ([#33](https://github.com/langchain-ai/langsmith-go/issues/33)) ([675dfba](https://github.com/langchain-ai/langsmith-go/commit/675dfba4834e9de34b55580f9cb492061acb7698))
* **deps:** bump google.golang.org/grpc ([#37](https://github.com/langchain-ai/langsmith-go/issues/37)) ([bd75dd5](https://github.com/langchain-ai/langsmith-go/commit/bd75dd510aed9e68fbb4cc1ee0d7ca0673a847a1))
* **internal:** codegen related update ([8beb884](https://github.com/langchain-ai/langsmith-go/commit/8beb884529d7bc3dde4ed9cbeb670350f46c6832))
* **internal:** codegen related update ([de31cba](https://github.com/langchain-ai/langsmith-go/commit/de31cba21aca4195cdc3be341e4a5882d3d8da53))
* **internal:** minor cleanup ([3bd3438](https://github.com/langchain-ai/langsmith-go/commit/3bd3438e0341b55894667b2979b9e80c1ee788f9))
* **internal:** move custom custom `json` tags to `api` ([dac02d6](https://github.com/langchain-ai/langsmith-go/commit/dac02d65bd10b3c3b0978ef881cf78cb7c1f49c1))
* **internal:** remove mock server code ([29de91e](https://github.com/langchain-ai/langsmith-go/commit/29de91e1d14eb45ec8149e451f6b3909cbecb9f9))
* **internal:** tweak CI branches ([bb0d1e5](https://github.com/langchain-ai/langsmith-go/commit/bb0d1e50b42591f400aad576ae04d030a723229d))
* **internal:** use explicit returns ([a45d5f9](https://github.com/langchain-ai/langsmith-go/commit/a45d5f9d16e30714eea65fa1c853202f0b509ffc))
* **internal:** use explicit returns in more places ([49ac557](https://github.com/langchain-ai/langsmith-go/commit/49ac557268964fee91b086ca93d26941bcd4bdcd))
* update mock server docs ([08bd0c3](https://github.com/langchain-ai/langsmith-go/commit/08bd0c3bbe4a380ef361cbedcef438996806a78b))
* update placeholder string ([5165a27](https://github.com/langchain-ai/langsmith-go/commit/5165a27b225207cbfd3c9f1aed806a4e30e0db5d))

## 0.1.0 (2026-02-19)

Full Changelog: [v0.1.0-alpha.10...v0.1.0](https://github.com/langchain-ai/langsmith-go/compare/v0.1.0-alpha.10...v0.1.0)

### Chores

* stable release ([#28](https://github.com/langchain-ai/langsmith-go/issues/28)) ([127e4a7](https://github.com/langchain-ai/langsmith-go/commit/127e4a7027ff8066289a4f0407a07730209e2ee0))
* stable release ([#31](https://github.com/langchain-ai/langsmith-go/issues/31)) ([3e69fc6](https://github.com/langchain-ai/langsmith-go/commit/3e69fc60056c409c53389d5ccb2746a93eaf0107))

## 0.1.0-alpha.10 (2026-02-19)

Full Changelog: [v0.1.0-alpha.9...v0.1.0-alpha.10](https://github.com/langchain-ai/langsmith-go/compare/v0.1.0-alpha.9...v0.1.0-alpha.10)

### Features

* **api:** api update ([fafcd6d](https://github.com/langchain-ai/langsmith-go/commit/fafcd6dac4bb2999e9caf0ecd48e160814d4f20a))
* **api:** api update ([a37c8c2](https://github.com/langchain-ai/langsmith-go/commit/a37c8c2dc4938ea718370bd1a1cfc83e63ffd852))
* **api:** api update ([e892367](https://github.com/langchain-ai/langsmith-go/commit/e89236788b500e796c37d6f19150acc0e96444af))
* **api:** api update ([b0f5add](https://github.com/langchain-ai/langsmith-go/commit/b0f5addf212f54d27f50700b8938317e04477030))
* **api:** api update ([3baff5d](https://github.com/langchain-ai/langsmith-go/commit/3baff5dcc7247fd4191e91715aab109c0d644536))
* **api:** api update ([4bb8835](https://github.com/langchain-ai/langsmith-go/commit/4bb88355dc1215542f0afdfd21919f38c3b81b52))
* **api:** api update ([79854d5](https://github.com/langchain-ai/langsmith-go/commit/79854d5907831d0ccbfd611062ccf03d0baf22ed))
* **api:** api update ([488b81d](https://github.com/langchain-ai/langsmith-go/commit/488b81d21399c46ab5338c7e23552e3e790564a3))
* **api:** api update ([735e948](https://github.com/langchain-ai/langsmith-go/commit/735e9482c7cfdb19554688a08a94055572d104e1))


### Bug Fixes

* **client:** use correct format specifier for header serialization ([7f9dfea](https://github.com/langchain-ai/langsmith-go/commit/7f9dfeafa7d1272a77b681c6d04f4e802577ac94))


### Chores

* **api:** minor updates ([d251895](https://github.com/langchain-ai/langsmith-go/commit/d2518950220c7132bfeabd5fe37db49d4f2ae34f))

## 0.1.0-alpha.9 (2026-01-27)

Full Changelog: [v0.1.0-alpha.8...v0.1.0-alpha.9](https://github.com/langchain-ai/langsmith-go/compare/v0.1.0-alpha.8...v0.1.0-alpha.9)

### Features

* **api:** api update ([3a33ebc](https://github.com/langchain-ai/langsmith-go/commit/3a33ebce78a90e57932eb790087e81f9bb85dbbe))
* **api:** api update ([276df7b](https://github.com/langchain-ai/langsmith-go/commit/276df7b1a02c02699b3e846896ac02fe3c9ea1ed))
* **api:** manual updates ([759e660](https://github.com/langchain-ai/langsmith-go/commit/759e6604655570b6387f89ab01d3a70496b95116))
* **api:** manual updates ([fb96e0d](https://github.com/langchain-ai/langsmith-go/commit/fb96e0daa29ceced88eff17b6bcce8a2a1b55966))

## 0.1.0-alpha.8 (2026-01-23)

Full Changelog: [v0.1.0-alpha.7...v0.1.0-alpha.8](https://github.com/langchain-ai/langsmith-go/compare/v0.1.0-alpha.7...v0.1.0-alpha.8)

### Features

* **api:** api update ([086dbef](https://github.com/langchain-ai/langsmith-go/commit/086dbefc54d0aa66112df736ffa9931db9023475))
* **api:** manual updates ([6df9292](https://github.com/langchain-ai/langsmith-go/commit/6df92927c7dc7072023b115552942b8641b94b65))
* **api:** manual updates ([7146ed4](https://github.com/langchain-ai/langsmith-go/commit/7146ed43004342d207bbd6250349ae2f77f62f5d))
* **api:** manual updates ([531ec16](https://github.com/langchain-ai/langsmith-go/commit/531ec16a61dba409170587de365832a09e3326c6))


### Chores

* rename langchain_base_url to langsmith_endpoint ([#22](https://github.com/langchain-ai/langsmith-go/issues/22)) ([66b5cb4](https://github.com/langchain-ai/langsmith-go/commit/66b5cb449050d5fff16ed86bad49fe82a99b3eee))

## 0.1.0-alpha.7 (2026-01-22)

Full Changelog: [v0.1.0-alpha.6...v0.1.0-alpha.7](https://github.com/langchain-ai/langsmith-go/compare/v0.1.0-alpha.6...v0.1.0-alpha.7)

### Features

* added high level tracer initialization ([#21](https://github.com/langchain-ai/langsmith-go/issues/21)) ([031e9e8](https://github.com/langchain-ai/langsmith-go/commit/031e9e89a644ca7e5ca518e3c8bbec50461e0585))
* **api:** api update ([146db7b](https://github.com/langchain-ai/langsmith-go/commit/146db7b1b43663bc4d064aa8b40c4faff8339942))


### Bug Fixes

* **docs:** add missing pointer prefix to api.md return types ([0d8ba3a](https://github.com/langchain-ai/langsmith-go/commit/0d8ba3ae02d75b410a89cee54d07f7e7b73ba7ea))


### Chores

* **internal:** update `actions/checkout` version ([faac99e](https://github.com/langchain-ai/langsmith-go/commit/faac99e1ba11d1b495e077bfcc84f7b3336e280c))

## 0.1.0-alpha.6 (2026-01-15)

Full Changelog: [v0.1.0-alpha.5...v0.1.0-alpha.6](https://github.com/langchain-ai/langsmith-go/compare/v0.1.0-alpha.5...v0.1.0-alpha.6)

### Features

* **api:** api update ([402c765](https://github.com/langchain-ai/langsmith-go/commit/402c7650086efacfee5d42d11bf5672745da8f55))


### Bug Fixes

* added fix for missing thread token counts ([#19](https://github.com/langchain-ai/langsmith-go/issues/19)) ([6c54b0f](https://github.com/langchain-ai/langsmith-go/commit/6c54b0fbac81a2df9b0de9d38f3e3fd54180637d))

## 0.1.0-alpha.5 (2026-01-14)

Full Changelog: [v0.1.0-alpha.4...v0.1.0-alpha.5](https://github.com/langchain-ai/langsmith-go/compare/v0.1.0-alpha.4...v0.1.0-alpha.5)

### Features

* **api:** api update ([925146c](https://github.com/langchain-ai/langsmith-go/commit/925146c93683a7b98324d8b1b2bf27b3d4deec8d))
* **api:** api update ([40886b9](https://github.com/langchain-ai/langsmith-go/commit/40886b91e1471e64dff8dfb274b8bfe179ac99f9))
* **api:** api update ([74c71ed](https://github.com/langchain-ai/langsmith-go/commit/74c71eda0feb1f9a5885404d8583671e34585fb1))
* **api:** api update ([9dad560](https://github.com/langchain-ai/langsmith-go/commit/9dad56055048b0a1376c45d125de51a9da8f7446))
* **api:** api update ([df983a3](https://github.com/langchain-ai/langsmith-go/commit/df983a31e6f0ded956237952f91aade41bed48aa))
* **api:** api update ([13e585b](https://github.com/langchain-ai/langsmith-go/commit/13e585bda50d281da4a45d9968acdf37567555cf))
* **api:** api update ([e7ba714](https://github.com/langchain-ai/langsmith-go/commit/e7ba714691bc74fb530a729ece4f1bf63159d2a4))
* **api:** api update ([64f64ed](https://github.com/langchain-ai/langsmith-go/commit/64f64ed425d061c542cdc767c2991c3065b28f30))
* **api:** manual updates ([13e2d39](https://github.com/langchain-ai/langsmith-go/commit/13e2d39dd5518d343d1dbd1fd5ef07744f9e0c9a))


### Bug Fixes

* skip usage tests that don't work with Prism ([6f397c4](https://github.com/langchain-ai/langsmith-go/commit/6f397c4441cd81e2d6ff8e85953502e1e39a8afb))


### Chores

* **internal:** codegen related update ([cbcbc79](https://github.com/langchain-ai/langsmith-go/commit/cbcbc790d92feb7a25fce8dcccdf363504e76a96))

## 0.1.0-alpha.4 (2025-12-10)

Full Changelog: [v0.1.0-alpha.3...v0.1.0-alpha.4](https://github.com/langchain-ai/langsmith-go/compare/v0.1.0-alpha.3...v0.1.0-alpha.4)

### Features

* added tracing with go openai client ([#16](https://github.com/langchain-ai/langsmith-go/issues/16)) ([49d5dce](https://github.com/langchain-ai/langsmith-go/commit/49d5dce40aced07d8fef90dd4bcc5eb6391221fc))
* **api:** api update ([5a2210d](https://github.com/langchain-ai/langsmith-go/commit/5a2210dafdd9faa25832178a5edc92d92b9dac58))

## 0.1.0-alpha.3 (2025-12-09)

Full Changelog: [v0.1.0-alpha.2...v0.1.0-alpha.3](https://github.com/langchain-ai/langsmith-go/compare/v0.1.0-alpha.2...v0.1.0-alpha.3)

### Features

* **api:** api update ([3b81dd1](https://github.com/langchain-ai/langsmith-go/commit/3b81dd1e9e37c6d7c3932dce8ce8730fcd9b7515))
* go-openai and anthropic tracing examples ([#14](https://github.com/langchain-ai/langsmith-go/issues/14)) ([cdddc56](https://github.com/langchain-ai/langsmith-go/commit/cdddc56d3d9c6bfa7a7be385757a2b05b9ebe139))

## 0.1.0-alpha.2 (2025-12-09)

Full Changelog: [v0.1.0-alpha.1...v0.1.0-alpha.2](https://github.com/langchain-ai/langsmith-go/compare/v0.1.0-alpha.1...v0.1.0-alpha.2)

### Features

* **api:** api update ([9961cd7](https://github.com/langchain-ai/langsmith-go/commit/9961cd783e6a42711329554a16ffe8a4ecdbe175))
* **api:** manual updates ([5b09d69](https://github.com/langchain-ai/langsmith-go/commit/5b09d69af8b36c0b824ca47c9c7b48ef7f9e9339))


### Chores

* **api:** delete index methods ([a9333cc](https://github.com/langchain-ai/langsmith-go/commit/a9333cc11eed5b3061e21317e5f0b2ee06b08041))

## 0.1.0-alpha.1 (2025-12-09)

Full Changelog: [v0.0.1...v0.1.0-alpha.1](https://github.com/langchain-ai/langsmith-go/compare/v0.0.1...v0.1.0-alpha.1)

### Features

* add prompt example ([#2](https://github.com/langchain-ai/langsmith-go/issues/2)) ([f9d7cea](https://github.com/langchain-ai/langsmith-go/commit/f9d7ceaedff9477b12b95e768352fdf2e400a132))
* add record experiment example ([#3](https://github.com/langchain-ai/langsmith-go/issues/3)) ([36834e3](https://github.com/langchain-ai/langsmith-go/commit/36834e34da8a91ff1e03e2c7397b2c277e26bd2e))
* added dataset example ([#7](https://github.com/langchain-ai/langsmith-go/issues/7)) ([a575630](https://github.com/langchain-ai/langsmith-go/commit/a57563027011e78d36f484e4f1c0dafa18d2fc5c))
* added evaluation example ([#8](https://github.com/langchain-ai/langsmith-go/issues/8)) ([df85f8f](https://github.com/langchain-ai/langsmith-go/commit/df85f8f0ea904aa0152f67af33915ae9dfcb6c14))
* added otel ingestion examples ([de957e9](https://github.com/langchain-ai/langsmith-go/commit/de957e93220576ce553ace2b2de0bc5d308e1108))
* **api:** add go target ([07c2406](https://github.com/langchain-ai/langsmith-go/commit/07c240666a7c39e1d684090b20286dd3e5dae6ec))
* **api:** api update ([720fae2](https://github.com/langchain-ai/langsmith-go/commit/720fae2738cf1101e291dd1caf849bda2030379b))
* **api:** api update ([f55b0d3](https://github.com/langchain-ai/langsmith-go/commit/f55b0d30a81c1b0b9728b00c30875452de31e8aa))
* **api:** api update ([afc0dd8](https://github.com/langchain-ai/langsmith-go/commit/afc0dd80d6a758e51ada8b308d6fc65bf9fc4d6b))
* **api:** api update ([92faff9](https://github.com/langchain-ai/langsmith-go/commit/92faff9c16d4f7a1026d9dd9c1a85fe37afa2284))
* **api:** api update ([4620906](https://github.com/langchain-ai/langsmith-go/commit/4620906243fa4a98071c89dcef2e815b13e5e626))
* **api:** api update ([5800689](https://github.com/langchain-ai/langsmith-go/commit/5800689a7207e62f058ec2ba8062f4d955496855))
* **api:** manual updates ([e360eb9](https://github.com/langchain-ai/langsmith-go/commit/e360eb9ad3f42612329f3b5e5ad397317a1f4b0a))
* **api:** pagination read me ([1f39895](https://github.com/langchain-ai/langsmith-go/commit/1f398959ad0828412ff7a25cca1511f18899a714))
* **api:** readme example - change REPLACE_ME session_id to uuid ([b4014e5](https://github.com/langchain-ai/langsmith-go/commit/b4014e5653f306ab5da7221d0ac482e8cbd08911))
* list runs example ([218ea3b](https://github.com/langchain-ai/langsmith-go/commit/218ea3b5914409190927e130699e0d2ea9c4ffc9))


### Bug Fixes

* **mcp:** correct code tool API endpoint ([2088a47](https://github.com/langchain-ai/langsmith-go/commit/2088a47e50a2dabf924637810e0d1ff53a29b096))
* prompt example dotted order formatting ([#10](https://github.com/langchain-ai/langsmith-go/issues/10)) ([8961794](https://github.com/langchain-ai/langsmith-go/commit/8961794aaf7f488031ec95581e403101ba742e70))
* rename param to avoid collision ([0b45bf0](https://github.com/langchain-ai/langsmith-go/commit/0b45bf08560d5661d52377a622d873f393865f05))


### Chores

* elide duplicate aliases ([4bd5c33](https://github.com/langchain-ai/langsmith-go/commit/4bd5c33d7ee3f8b17439449b5fadf510f949181b))
* **internal:** codegen related update ([d98b055](https://github.com/langchain-ai/langsmith-go/commit/d98b0552999fd1315767a9d0a1232a40d066d5c3))
* update SDK settings ([d6cb6f1](https://github.com/langchain-ai/langsmith-go/commit/d6cb6f18a01dedce7b1ec43c3ab59a8cf1915ade))


### Documentation

* Add configuration to README ([#5](https://github.com/langchain-ai/langsmith-go/issues/5)) ([59244be](https://github.com/langchain-ai/langsmith-go/commit/59244bec89a9804cfac62f1e7a550b6f5fa48329))
* add exmaple section to README ([#9](https://github.com/langchain-ai/langsmith-go/issues/9)) ([70046c7](https://github.com/langchain-ai/langsmith-go/commit/70046c7b0f6a6cdf115a9d6adf28ee62920868c7))
