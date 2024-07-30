from typing import List

from pydantic import BaseModel
from sqlalchemy import Column, Integer, String, create_engine
from sqlalchemy.orm import Session, declarative_base, sessionmaker

from fastapi import Depends, FastAPI, HTTPException

# FastAPI app instance
app = FastAPI()

# Database setup
DATABASE_URL = "sqlite:///../../items.db"
engine = create_engine(DATABASE_URL)
SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)
Base = declarative_base()


# Database model
class Item(Base):
    __tablename__ = "items"
    id = Column(Integer, primary_key=True, index=True)
    name = Column(String, index=True)
    description = Column(String)


# Create tables
Base.metadata.create_all(bind=engine)


# Dependency to get the database session
def get_db():
    db = SessionLocal()
    try:
        yield db
    finally:
        db.close()


# Pydantic model for request data
class ItemCreate(BaseModel):
    name: str
    description: str


# Pydantic model for response data
class ItemResponse(BaseModel):
    id: int
    name: str
    description: str


@app.post("/items/", response_model=ItemResponse)
async def create_item(item: ItemCreate, db: Session = Depends(get_db)) -> ItemResponse:
    """API endpoint to create an item"""
    db_item = Item(**item.dict())
    db.add(db_item)
    db.commit()
    db.refresh(db_item)
    return db_item


@app.get("/items/")
async def list_items(db: Session = Depends(get_db)) -> List[ItemResponse]:
    """API endpoint to list all Items"""
    return db.query(Item).all()


@app.get("/items/{item_id}", response_model=ItemResponse)
async def read_item(item_id: int, db: Session = Depends(get_db)) -> ItemResponse:
    """API endpoint to read an item by ID"""
    db_item = db.query(Item).filter(Item.id == item_id).first()
    if db_item is None:
        raise HTTPException(status_code=404, detail="Item not found")
    return db_item


if __name__ == "__main__":
    import uvicorn

    # Run the FastAPI application using Uvicorn
    uvicorn.run(app, host="127.0.0.1", port=8000)
