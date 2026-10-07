# TASK-010-T3-FT-005-W4 — independent acceptance evidence

FT-005-AC-002 — Два места удаления и текущие условия

Observable criterion: удаление из списка и настроек убирает выбранный поиск и его историю, сохраняя аренду и остальные поиски; старое меню не восстанавливает удалённый ID. Перед новой доставкой проверяется существование/актуальность поиска; устаревшая оценка после правки не становится отказом для новых условий. Уже отправленный Telegram запрос не отзывается.

Fresh verifier proof: TestVerifierDeleteMenuOutcome (list-private/settings-target), TestVerifierDeleteRollbackExactRows, TestVerifierStaleEvaluationOutcome (delete/edit-reject/edit-send), TestVerifierStartedRequestDeletionOutcome. Race probe exit0 covers every mapped current harm-driving result; focused native package tests and scripts/build.sh repeated exit0. Full mapping, architecture review, executor RED/GREEN assessment and T3 isolation in `.protocols/TASK-010-T3-FT-005-W4/verification.md`.

Reproduce: `sh .tasks/TASK-010-T3-FT-005-W4/verifier-probes.sh`. Command results and gate stdin reproduction: `verifier-commands-results.json`. Functional trace: `verifier-outcomes.log`; focused/native logs, pre/final source hashes, fresh authored test artifacts share this folder. Initial shell setup error is recorded but excluded from functional evidence. No execute receipt reuse. Source/HEAD unchanged; co-review best-effort Sol xhigh unavailable thread limit, no retry. Lifecycle unchanged for root; per-task T3 red-verify due.
