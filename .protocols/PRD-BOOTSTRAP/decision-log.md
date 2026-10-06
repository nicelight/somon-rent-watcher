# PRD bootstrap decision log

- Existing `seen_ads` semantics remain authoritative; no shown-history feature was created.
- Fallback selection and Telegram interaction are separate features because each has its own user-visible result and verification boundary.
- Production release is a separate feature because it has an independently blocking outcome, operator authorization, runtime boundary, and rollback/isolation acceptance.
- No unresolved decomposition decision remains.

## Accepted keyword-monitoring delta — 2026-10-06

- Источник: принятый раздел `.memory-bank/prd.md#keyword-monitoring-proposal--2026-10-06`, включая подтверждение оператора «Да, эти правила подходят».
- EP-002 / FT-005 — один самостоятельный пользовательский результат: сохраняемый поиск с собственными условиями, lifecycle и уведомлениями. Настройка/удаление являются его lifecycle; storage/UI/source не разделены на продуктовые features.
- Продуктовые решения закрыты; scope-binding — пробел source evidence, а не повод менять acceptance. Его сохраняет FT-005 SDD gate перед source-dependent task handoff/исполнением.
- FT-001…004, существующие ID, задачи/approval и Planning Revision 1 не изменяются.
