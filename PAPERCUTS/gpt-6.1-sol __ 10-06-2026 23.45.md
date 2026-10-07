# Наблюдения текущей сессии

- Existing `somon.ParseDetail` принимает произвольный неполный HTML, если fallback card уже содержит ID/title. Это подтверждено локальной detail fixture `broken detail`; такой body не является воспроизводимой parse error этого adapter. В TASK-009 source parse failure проверяется реальной ошибкой keyword выдачи, detail retry — HTTP 500. Исправление parser вне этой задачи.
