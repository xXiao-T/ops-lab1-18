# 实验 1 记录：第 18 组

现场检查时将打开本文件，按其中的链接核对。证据一律以链接或文本形式给出，不要只贴截图。将 `<填写>` 替换为本组的内容。main 受保护之后，修改本文件也须通过 PR。

## 1. 基本信息

| 项 | 内容 |
|---|---|
| 小组 | 第 18 组 |
| 仓库地址 | <https://github.com/xXiao-T/ops-lab1-18> |
| 代码载体 | 统一示例（enroll-conflict，Go） |
| runner 与语言版本 | `ubuntu-24.04`，Go 1.26.9 |
| 合并方式 | squash（仅允许 squash 合并 + 合并后自动删除分支） |

## 2. 任务 0：本地基线

| 问题 | 观察结果 |
|---|---|
| 第一次 `make package` 得到的 BUILD_INFO（粘贴） | 见下方代码块 |
| 修改一行、不提交，再次打包：BUILD_INFO 是否变化？二进制文件的 sha256 是否变化？（贴出两次的 sha256） | BUILD_INFO **不变**（`commit` 取自 `git rev-parse HEAD`，不含工作区状态）；二进制 sha256 **变化**：`f23447d406c42af7dd35e23c2d1b2f15484cfb709cc9b4bb1262643694b8528e` → `9778e72f52340f28797756e9cb7df8889111c5f7dac452a737411f1426e2911c` |
| 两次 `go version -m` 输出中的 vcs 行有何不同？ | `vcs.modified` 由 `false` 变为 `true`；`vcs.revision`（`7b63b41c…`）与 `vcs.time` 不变 |
| 仅凭 BUILD_INFO，能否回答该包是什么、在哪台机器上、何时、用哪份代码构建？哪一项信息会出错？ | 「是什么」能答（`name`）；「何时」「在哪台机器」**不能**——BUILD_INFO 有意不含构建时间与主机信息（保证可重复构建）。「用哪份代码」**会出错**：源码已改、`commit` 仍指向 `7b63b41c`，与实际构建内容不符（只有 `vcs.modified=true` 暴露了这一点）。故本地手工构建的标识不可信，需要流水线以 `github.sha` 等权威输入生成标识 |

第一次的 BUILD_INFO：

```text
name=enroll-conflict
version=1.0.0
build_no=1
run_id=local
commit=7b63b41c13f061439997b0e6e4eaa2af922a7aa1
go=go1.26.9
```

## 3. 证据

| 项 | 链接或文本 |
|---|---|
| 规则集里启用的规则 | `deletion`（Restrict deletions）、`non_fast_forward`（Block force pushes）、`pull_request`（Required approvals = **1**）、`required_status_checks`（必需检查 = **`ci`**）；Bypass list 为空（对管理员同样生效） |
| 组长直接 `git push origin main` 被拒的终端输出 | 见下方代码块 |
| 第一轮：失败的 PR 运行 | https://github.com/xXiao-T/ops-lab1-18/actions/runs/38859655635 （PR #1，倒在「运行测试」） |
| 第一轮：修复后通过的 PR 运行 | https://github.com/xXiao-T/ops-lab1-18/actions/runs/38859793662 |
| 第一轮：合入 main 后的运行（含制品） | https://github.com/xXiao-T/ops-lab1-18/actions/runs/39021031583 （制品 `enroll-conflict-8177aa7…-4`） |
| 第二轮：找到的缺陷，以及现有测试为何无法发现它 | `Valid()` 把 `Weekday > 7` 收紧为 `> 6`，周日（7）被误判非法；现有测试的 Weekday 只覆盖 0、3、8（单元）与 1、9（端点），全仓库没有 `Weekday=7` 的用例，故测试全绿 |
| 第二轮：仅含缺陷时仍然通过的运行 | https://github.com/xXiao-T/ops-lab1-18/actions/runs/39021467507 （gate 未拦住）|
| 第二轮：补测试后失败的运行 | https://github.com/xXiao-T/ops-lab1-18/actions/runs/39021583380 （新用例抓获）|
| 第二轮：撤销缺陷后通过的运行 | https://github.com/xXiao-T/ops-lab1-18/actions/runs/39101050460 |
| 未合入分支上的 tag 被 release gate 拒绝的运行 | https://github.com/xXiao-T/ops-lab1-18/actions/runs/39121326245 （`release gate（tag 必须在 main 上）` 失败，后续步骤全部跳过，未产生 Release）|
| Release v1.0.0 | https://github.com/xXiao-T/ops-lab1-18/releases/tag/v1.0.0 |
| Release v1.0.1 | <https://github.com/xXiao-T/ops-lab1-18/releases/tag/v1.0.1> |

组长直接 `git push origin main` 被拒的输出：

```text
remote: error: GH013: Repository rule violations found for refs/heads/main.
remote: - Changes must be made through a pull request.
remote: - Required status check "ci" is expected.
 ! [remote rejected] main -> main (push declined due to repository rule violations)
error: failed to push some refs to 'https://github.com/xXiao-T/ops-lab1-18.git'
```

## 4. 追溯与运行

对 v1.0.0 和 v1.0.1 各运行一次 `scripts/trace.sh`，粘贴输出，并填写下表。

| 项 | v1.0.0 | v1.0.1 |
|---|---|---|
| 版本号 | `1.0.0` | `1.0.1` |
| 构建号 | `1` | `3` |
| 提交号（SHA） | `8a003ec406d0946494f606064e5ae6e86c51d005` | `e19fc9c92a09011af2a4331887adcb1e936868e7` |
| 引入该提交的 PR | #3 <https://github.com/xXiao-T/ops-lab1-18/pull/3> | #4 <https://github.com/xXiao-T/ops-lab1-18/pull/4> |
| 构建它的 Actions 运行（`trace.sh` 给出的链接） | <https://github.com/xXiao-T/ops-lab1-18/actions/runs/39114687305> | <https://github.com/xXiao-T/ops-lab1-18/actions/runs/39123318787> |
| 制品 sha256（发布当时，取自 Release 里的 .sha256） | `3fca5242d4a042baaf7d4384b60f9d9b39708d7813dba75daa6b02830703a949` | `5ac2510fa2359664079193bbcb42fec09b6415dc2ad24f3dab0da30bb3d533c1` |

**仅凭构建号能否唯一确定是哪次运行？** 不能。构建号在每个工作流内各自从 1 计数：v1.0.0 的构建号 1 在 `release` 工作流中对应第 1 次运行（<https://github.com/xXiao-T/ops-lab1-18/actions/runs/39114687305>），而 `ci` 工作流也有第 1 次运行（<https://github.com/xXiao-T/ops-lab1-18/actions/runs/38846220491>，main 首次 push）。两者同号，仅凭构建号无法区分是哪个工作流的哪次运行；`run_id` 全局唯一，才能唯一定位。

**发布 v1.0.1 之后重新下载 v1.0.0 的附件**，计算得到的 sha256：`3fca5242d4a042baaf7d4384b60f9d9b39708d7813dba75daa6b02830703a949`。与上表中发布当时的值**一致**——Release 制品不可变。

`trace.sh` 输出（v1.0.0）：

```text
== 1. 校验和
enroll-conflict_1.0.0_linux_amd64.tar.gz: OK
== 2. 制品里的 BUILD_INFO
name=enroll-conflict
version=1.0.0
build_no=1
run_id=39114687305
commit=8a003ec406d0946494f606064e5ae6e86c51d005
go=go1.26.9
== 3. 提交 8a003ec406d0946494f606064e5ae6e86c51d005
作者: xXiao-T <1545535109@sjtu.edu.cn>
时间: 2026-10-18T17:10:18+08:00
标题: ci: 添加 release 工作流（release gate + 发布构建 + Release） (#3)
主干: 该提交在 origin/main 的历史上
指向它的 tag: v1.0.0
版本号: 1.0.0    构建号: 1
提交页面（含关联的 PR）: https://github.com/xXiao-T/ops-lab1-18/commit/8a003ec406d0946494f606064e5ae6e86c51d005
构建它的运行: https://github.com/xXiao-T/ops-lab1-18/actions/runs/39114687305
```

`trace.sh` 输出（v1.0.1）：

```text
== 1. 校验和
enroll-conflict_1.0.1_linux_amd64.tar.gz: OK
== 2. 制品里的 BUILD_INFO
name=enroll-conflict
version=1.0.1
build_no=3
run_id=39123318787
commit=e19fc9c92a09011af2a4331887adcb1e936868e7
go=go1.26.9
== 3. 提交 e19fc9c92a09011af2a4331887adcb1e936868e7
作者: xXiao-T <1545535109@sjtu.edu.cn>
时间: 2026-10-18T19:27:30+08:00
标题: docs: 填写 LAB1.md 第 1–3 节（基本信息、任务 0、证据） (#4)
主干: 该提交在 origin/main 的历史上
指向它的 tag: v1.0.1
版本号: 1.0.1    构建号: 3
提交页面（含关联的 PR）: https://github.com/xXiao-T/ops-lab1-18/commit/e19fc9c92a09011af2a4331887adcb1e936868e7
构建它的运行: https://github.com/xXiao-T/ops-lab1-18/actions/runs/39123318787
```

任务 6 第 4 步，两次 `/version` 的输出：

```text
{"version":"1.0.1","commit":"e19fc9c92a09011af2a4331887adcb1e936868e7","build_no":"3","go":"go1.26.9"}
{"version":"1.0.0","commit":"8a003ec406d0946494f606064e5ae6e86c51d005","build_no":"1","go":"go1.26.9"}
```

## 5. 数据

时间均取自 GitHub（API/页面），时区为 **UTC**。

第一轮 gate 演练的时间线：

| 事件 | 时间 | 链接 |
|---|---|---|
| 发起 PR | 2026-10-16 06:09:47 | <https://github.com/xXiao-T/ops-lab1-18/pull/1> |
| 第一次 `ci` 失败（运行结束） | 2026-10-16 06:10:10 | <https://github.com/xXiao-T/ops-lab1-18/actions/runs/38859655635> |
| 推送修复提交 | 2026-10-16 ~06:11:15 | <https://github.com/xXiao-T/ops-lab1-18/actions/runs/38859793662>（该运行 06:11:19 开始） |
| `ci` 通过（运行结束） | 2026-10-16 06:11:33 | <https://github.com/xXiao-T/ops-lab1-18/actions/runs/38859793662> |
| 合入 main | 2026-10-17 09:24:52 | <https://github.com/xXiao-T/ops-lab1-18/pull/1> |

提交到发布的间隔（四个合入的 PR 全部记录）：

| PR | 第一个提交的时间 | 合入时间 | 对应 Release 的发布时间 | 提交到发布的间隔 |
|---|---|---|---|---|
| #1 | 2026-10-16 06:08:50 | 2026-10-17 09:24:52 | v1.0.0，2026-10-18 09:16:26 | 2 天 3 小时 7 分 |
| #2 | 2026-10-17 09:31:20 | 2026-10-18 07:29:30 | v1.0.1，2026-10-18 11:27:33 | 1 天 1 小时 56 分 |
| #3 | 2026-10-18 07:48:31 | 2026-10-18 09:10:19 | v1.0.0，2026-10-18 09:16:26 | 1 小时 27 分 |
| #4 | 2026-10-18 10:59:47 | 2026-10-18 11:27:30 | v1.0.1，2026-10-18 11:27:33 | 27 分 46 秒 |

`ci` 工作流的运行时长（全部列出，截至 2026-10-20 00:35 UTC，共 21 次）：

| 运行（run_number） | 运行 ID | 事件（PR / push） | 时长 | 结论 | 备注 |
|---|---|---|---|---|---|
| 第 1 次 | 38846220491 | push（main 首次） | 44 s | ✅ | 首次全量，含打包上传 |
| 第 2 次 | 38859655635 | PR #1 | 17 s | ❌ | 单双周误判被现有测试捕获 |
| 第 3 次 | 38859793662 | PR #1 | 14 s | ✅ | 修复后通过 |
| 第 4 次 | 39021031583 | push（main） | 15 s | ✅ | 含打包上传 |
| 第 5 次 | 39021467507 | PR #2 | 7 s | ✅ | 仅缺陷，「盲区」通过 |
| 第 6 次 | 39021583380 | PR #2 | 13 s | ❌ | 补测试后抓获 |
| 第 7 次 | 39101050460 | PR #2 | 10 s | ✅ | 撤销缺陷后通过 |
| 第 8 次 | 39107852003 | push（main） | 12 s | ✅ | 含打包上传 |
| 第 9 次 | 39109060163 | PR #3 | 14 s | ✅ | 通过 |
| 第 10 次 | 39114285373 | push（main） | 13 s → Re-run 8 s | ✅ | 含打包上传；两次 attempt 输出逐字节相同 |
| 第 11 次 | 39121442639 | PR #4 | 11 s | ✅ | 通过 |
| 第 12 次 | 39123312961 | push（main，PR #4 合入） | 14 s | ✅ | 含打包上传 |
| 第 13 次 | 39123821654 | PR #5 | 10 s | ✅ | 通过 |
| 第 14 次 | 39123848441 | PR #5 | 20 s | ✅ | 补充 v1.0.1 链接后 |
| 第 15 次 | 39123858728 | PR #5 | 10 s | ✅ | 组号标题修正后 |
| 第 16 次 | 39219123991 | PR #6（OO1220） | 12 s | ✅ | 通过 |
| 第 17 次 | 39223142631 | PR #7（peichow2） | 7 s | ❌ | gofmt 检查拦截（首版格式问题）|
| 第 18 次 | 39224353639 | PR #7 | 13 s | ✅ | 修复格式后通过 |
| 第 19 次 | 39347322088 | push（main，PR #7 合入） | 17 s | ✅ | 含打包上传 |
| 第 20 次 | 39347330326 | push（main，PR #6 合入） | 25 s | ✅ | 含打包上传 |
| 第 21 次 | 39348146617 | PR #5 | 12 s | ✅ | 冲突解决后通过 |

> 运行时长按 job 的 started_at → completed_at 计算，含 Actions 排队与初始化开销；共享 runner 负载会造成波动，样本量小。

## 6. 分析问题

每题先列出所用数据（可引用上面的表），再用一段话（4～8 句）写出结论，并说明结论的成立条件。

**问题 1：gate 的收益与成本。** 两轮演练中，gate 拦住了什么，没有拦住什么？用第 5 节的时间线计算：第一轮因失败多花了多少时间，其中等待 `ci` 的时间和人工修改的时间各占多少？一次 `ci` 大约需要多久，在“提交到发布的间隔”中占多大比例？第二轮补充测试之后，还有哪类问题是测试本身无法发现的？

**数据**：第 5 节时间线与 `ci` 时长表。第一轮失败运行 17 s、修复运行 14 s；从失败运行结束（06:10:10）到修复运行通过（06:11:33）约 83 s，其中等待 `ci` 约 31 s（17+14），人工定位与修改约 52 s。一次 `ci` 约 7～44 s，在 27 分～2 天 3 小时的"提交到发布间隔"中占比约 0.3%～2.7%。从修复通过（06:11:33）到审批（次日 09:19:20）约 27 小时。第二轮中，仅含缺陷时 gate 全绿放行；补上覆盖 `Weekday=7` 的测试后才被拦截。

**结论**：gate 能拦住"测试能表达的缺陷"，拦不住"测试没想到的输入与行为"——第一轮 17 秒内被现有测试抓住，第二轮 7 秒内放行了盲区缺陷。它的主要成本不是机器时间（单次 ci 秒级），而是人工环节：本组审批等待最长达约 27 小时，是第一轮失败修复成本的数百倍。因此测试质量决定 gate 的有效性，补测试本身就是提升门禁收益的手段。补充测试之后，测试本身仍无法发现需求理解错误（测试按错误理解编写）、并发/竞态、性能退化、依赖外部环境的缺陷，以及测试作者没有想到的输入组合。**成立条件**：样本仅两次演练与 4 个 PR，审批等待受人为节奏影响，不能外推为团队常态。

**问题 2：PR 中的构建与 main 上的构建。** 在 PR 触发的运行日志中，`github.sha` 与 `pull_request.head.sha` 是否相同？它们能否在 main 的历史中找到？答案与合并方式有关，请写明本组使用的方式。合并之后，main 上这次改动的提交号是哪个？据此说明，为什么待发布的制品只能从 main 上的提交或 tag 构建。

**数据**：PR #1 首次运行日志：`event_name=pull_request`、`ref=refs/pull/1/merge`、`github.sha=fac0dca7323bb899e68ffa711f90f80f6b2cb04c`、`pr_head_sha=8a490c6185ccccca5bea75485181293d4d751644`；修复运行中 `github.sha` 变为 `e37be8e96907cfb52555219cb3c9e54b40068aff`、`pr_head_sha=620c24fcce0f1c8946ca32c94feae79163382182`。合并后 main 上的提交为 squash 提交：`8177aa7`（#1）、`e040638`（#2）、`8a003ec`（#3）、`e19fc9c`（#4），标题末尾带 `(#编号)`。

**结论**：在 PR 触发的运行中，`github.sha` 是 GitHub 为 PR 合成的"预演合并提交"，随分支 head 每次变动重新生成，**不在 main 的历史中**；`pull_request.head.sha` 是分支头，经 squash 合并后也不会作为提交出现在 main 上。只有合并产生的 squash 提交（或直接打在主干提交上的 tag）才真正代表"进入主干的源码"。因此待发布制品只能从 main 上的提交或 tag 构建——否则 BUILD_INFO 中的 `commit` 将指向主干历史之外的提交，`trace.sh` 的"主干血统"检查会失败，追溯也就失去意义。**成立条件**：本组使用 squash 合并；若使用 merge commit，PR 的分支头会成为 main 的祖先、可在历史中找到，但 `github.sha`（合成合并提交）仍不在主干上；若使用 rebase 合并，提交号会被重写，结论同样指向"只从 main/tag 构建"。

**问题 3：可重复构建。** 在本地同一个提交上按下表构建，记录哈希值。④、⑤ 需临时去掉 Makefile 中 tar 的 `--mtime=@0`，完成后用 `git checkout -- Makefile` 恢复，不要提交。

| 构建 | 条件 | 二进制 sha256 | tar.gz sha256 |
|---|---|---|---|
| ① | `make package`，`VERSION`、`BUILD_NO`、`RUN_ID`、`COMMIT` 均指定固定值 | `681edb17580e1329f46de3128faf3fb90d8572b737b3a2e48bc930324874a54a` | `1265c9487f569a0c59c41093dd53c5c9dc12c19779bc433ec0ff9c965366284f` |
| ② | 与 ① 完全相同，再构建一次 | `681edb17…a54a`（同 ①） | `1265c948…284f`（同 ①） |
| ③ | 仅将 `BUILD_NO` 换成另一个值 | `e70cf601b3fc71edb96c07387e443e26cfd3aeef12000ecc012c27182b7143de` | `f74563cbdf9f1660851bb64c74751fb797ec18ce69d36c42f71614cd23fa8caf` |
| ④ | 去掉 `--mtime=@0`，其余同 ① | `131e942d5784f24476753c7279219067c24b33adbd61278de8f78993487c9601` | `406d31c39d5dc7c023aac66cf3b5a123473839b651fc162acc4c00c4dbf47547` |
| ⑤ | 同 ④，间隔几秒后再构建一次 | `131e942d…9601`（同 ④） | `c4d1e68417688199766fb88013e74289a581736d0efc114f242e7d5eed48d008` |
| Actions | 对 main 上的一次 push 运行（39114285373）执行 Re-run all jobs，比较两次尝试中 package 步骤打印的值 | 两次相同：`a36f31be8e1fcb6a68890b0a4265be51c76260bfe42984d2a57583a9ad8c4d2c` | 两次相同：`068417a30d7918ee8f703d4941ddc264ab61c44498a58f64545153763d4a4101` |

哪些相同、哪些不同，原因分别是什么？④ 的二进制文件与 ① 是否相同？如果不同，请结合任务 0 中看到的 vcs 信息分析。据此说明制品中应当包含什么、不应包含什么。另外，`release` 在同一个提交上再次编译，而没有直接发布 `ci` 产出的主干制品：为什么这样做？这与 L-05 第 11 页的“一次构建、多处部署”是否冲突，代价是什么？

**哪些相同、哪些不同，原因**：① 与 ② 完全相同——固定输入、干净工作区、`-trimpath`，加上 tar 固定排序/属主/`--mtime=@0` 与 `gzip -n`，构建是确定性的。③ 仅改 `BUILD_NO`，二进制与 tar 均变（buildNo 经 ldflags 写入，BUILD_INFO 也随之变化）。④ 与 ① 相比 tar 变（去掉 `--mtime=@0`，改用真实文件时间）；**二进制也变**——修改 Makefile 使工作区变脏，Go 把 `vcs.modified=true` 编入二进制，与任务 0 的观察完全一致。⑤ 与 ④ 相隔 3 秒：二进制相同（输入与脏状态相同），tar 不同（含秒级 mtime）。Actions Re-run 的两次尝试中，`BUILD_INFO` 与两个 sha256 完全相同——`run_number`、`run_id` 不因重跑而改变。

**结论**：制品中应包含版本号、构建号、提交号、`run_id` 等"来自输入"的确定性标识；不应包含构建时间戳、工作区脏状态等随构建时刻或环境变化的量，否则无法复现、无法比较。`release` 在 tag 提交上重新编译，是为了把版本号经 ldflags 写入二进制——这是本实验的简化，与"一次构建、多处部署"确有冲突：同一提交产生了 `ci` 主干制品与 Release 制品两份构建，若环境或工作区状态不同，二者可能不同；代价是制品不唯一，需以 Release 为唯一发布来源并用 sha256 核对，发布链路多一次构建。**成立条件**：本地与 runner 的 Go 同为 1.26.9、源码干净；若 Go 版本不同，或引入 cgo/外部依赖，可重复性需另行验证。

**问题 4：发布与回退。** 假设 v1.0.1 发布后发现缺陷，应当如何处理？为什么不应删除 v1.0.1 的 tag、再用同一版本号重新发布？任务 6 中将服务从 v1.0.1 切换回 v1.0.0 时，使用的是哪个文件，是否需要重新构建？请用两次 `/version` 的输出说明。如果改为从源码重新构建 v1.0.0，得到的文件与原附件是否相同（参考问题 3）？最后，如果有人删除 v1.0.0 的附件，换上一个新文件和对应的新 `.sha256`，`trace.sh` 能否发现？本组在第 4 节记录的值能否发现？可以查阅 GitHub 的 Immutable releases 设置的用途。

**数据**：第 4 节两次 `/version` 输出与 v1.0.0 附件的 sha256（发布当时与 v1.0.1 发布后重下载均为 `3fca5242d4a042baaf7d4384b60f9d9b39708d7813dba75daa6b02830703a949`）；问题 3 的构建矩阵。

**结论**：v1.0.1 发布后发现缺陷，应发布 v1.0.2 修订版并在 Release 说明中注明问题，而不是删除 tag 再用同一版本号重发——版本号是面向用户的兼容性承诺，制品不可变；删除或重发会使已下载制品的追溯断链，"哪个制品对应哪次构建"失去唯一答案。回退时直接使用 v1.0.0 的 Release 附件（`tar -xzf` 解包后的二进制），无需重新构建：两次 `/version` 显示版本、提交号、构建号随制品切换（`1.0.1`/`e19fc9c…`/`3` → `1.0.0`/`8a003ec…`/`1`）。若从源码重建 v1.0.0，在固定 Go 版本、干净工作区且 BuildInfo 相同的前提下二进制可复现，但 tar.gz 受 mtime 与构建环境影响，不保证与原附件一致，因此回退应使用原附件。`trace.sh` 只验证"附件与同目录 `.sha256` 一致"；若有人把附件和 `.sha256` 一起替换，`trace.sh` 无法发现——但第 4 节记录的发布当时 sha256（`3fca5242…`）可以核出不一致；GitHub 的 Immutable releases 设置可在平台层面锁定已发布附件，防止替换。**成立条件**：第 4 节记录的值可信且保存在受保护的仓库历史中。

## 7. 选做记录

| 选做 | 完成情况与证据 |
|---|---|
| O1 容器镜像 | 未做 |
| O2 缓存对反馈时间的影响 | 未做 |
| O3 迁移到小组项目 | 未做 |

## 8. 成员贡献

每位成员一行，只填写 GitHub 用户名，无需填写姓名；行数应与小组人数一致。

| GitHub 用户名 | 发起并合入的 PR 编号 | 审批过的 PR 编号 |
|---|---|---|
| xXiao-T | #1、#2、#3、#4、#5 | #6、#7 |
| ZhenyangNi | [#8](https://github.com/xXiao-T/ops-lab1-18/pull/8)（已发起，待审批合入） | #1、#2、#3、#4 |
| OO1220 | #6 | 待 |
| peichow2 | #7 | 待 |

> 备注：`xXiao-T` 发起的 PR #5 已合入。`ZhenyangNi` 在 PR #8 中补充了 `/api/conflict` 的 13 个参数校验场景：分别验证 A、B 的星期上下界、起始节次、节次顺序及未知/空周次，并验证周日单节课合法且冲突。PR #8 须通过 `ci`，经另一位组员审批并 squash 合入后，才满足“每位成员至少有一个由本人发起并合入的 PR”的硬性要求。
