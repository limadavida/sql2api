from fastapi import FastAPI
from typing import List
from controllers import fetch_todos, create_todo, fetch_todo, update_todo, delete_todo
from schemas import TodoCreate, TodoResponse
from models import database

app = FastAPI()


@app.on_event(event_type="startup")
async def startup() -> None:
    await database.connect()


@app.on_event(event_type="shutdown")
async def shutdown() -> None:
    await database.disconnect()


#: Create
@app.post(path="/todos", response_model=TodoResponse)
async def add_todo(todo: TodoCreate) -> TodoResponse:
    return await create_todo(todo=todo)


#: Read All
@app.get(path="/todos", response_model=List[TodoResponse])
async def read_todos() -> List[TodoResponse]:
    return await fetch_todos()


#: Read One
@app.get(path="/todos/{todo_id}", response_model=TodoResponse)
async def read_todo(todo_id: int) -> TodoResponse:
    return await fetch_todo(todo_id=todo_id)


#: Update One
@app.put(path="/todos/{todo_id}", response_model=TodoResponse)
async def edit_todo(todo_id: int, todo: TodoCreate) -> TodoResponse:
    return await update_todo(todo_id, todo)


#: Delete One
@app.delete(path="/todos/{todo_id}")
async def del_todo(todo_id: int) -> None:
    await delete_todo(todo_id)


if __name__ == "__main__":
    import uvicorn

    uvicorn.run(app=app, host="127.0.0.1", port=8000)
