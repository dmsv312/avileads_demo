# AviLeads Web

## Добавление роли
1. `models.rule.go` добавить в _const_ новую роль
2. `models.rule.go` добавить в _rule_ новую роль c описание, так она будет появлятся в списках ролей в UI.
3. Если это новый пункт меню.  
3.1 Добавляем в `navbar.html`  
3.2 При необходимости добавляем в нужный `left-nav.html`  
4. Придумываем и прописывает `route` в файле `http/routers/router.go`
5. `http/routers/router.go` добавляем в `routeAccesses` с указанием роли из пункта `1`