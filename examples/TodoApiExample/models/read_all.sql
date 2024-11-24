-- name: get_all_tasks
-- description: usecase for task reading in context of TODO app
-- async: yes
-- parameters: max=10

SELECT * FROM todos;

