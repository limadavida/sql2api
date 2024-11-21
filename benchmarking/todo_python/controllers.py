from sqlalchemy.orm import Session
from fastapi import HTTPException
from models import TodoModel, database
from schemas import TodoCreate, TodoResponse


async def fetch_todos() -> list[TodoResponse]:
    query = "SELECT * FROM todos"
    todos = await database.fetch_all(query)
    return [
        TodoResponse(
            id=todo["id"],
            title=todo["title"],
            completed=(todo["completed"] if todo["completed"] is not None else False),
        )
        for todo in todos
    ]


async def create_todo(todo: TodoCreate) -> TodoResponse:
    query = TodoModel.__table__.insert().values(
        title=todo.title, completed=todo.completed
    )
    todo_id = await database.execute(query)
    return TodoResponse(id=todo_id, title=todo.title, completed=todo.completed)


async def fetch_todo(todo_id: int) -> TodoResponse:
    query = TodoModel.__table__.select().where(TodoModel.id == todo_id)
    todo = await database.fetch_one(query)

    if todo is None:
        raise HTTPException(status_code=404, detail="Todo not found")

    return TodoResponse(
        id=todo["id"],
        title=todo["title"],
        completed=(todo["completed"] if todo["completed"] is not None else False),
    )


async def update_todo(todo_id: int, todo_data: TodoCreate) -> TodoResponse:
    query = (
        TodoModel.__table__.update()
        .where(TodoModel.id == todo_id)
        .values(
            title=todo_data.title,
            completed=todo_data.completed,  # Não esqueça de incluir isso
        )
    )
    await database.execute(query)

    return await fetch_todo(todo_id)


async def delete_todo(todo_id: int) -> None:
    query = TodoModel.__table__.delete().where(TodoModel.id == todo_id)
    result = await database.execute(query)

    if result == 0:
        raise HTTPException(status_code=404, detail="Todo not found")
    return {"status": "OK"}
