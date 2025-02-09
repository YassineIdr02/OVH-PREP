package workflows

import (
	"fmt"
	"ovh-prep/models" 
	"go.temporal.io/sdk/workflow"
)

type ReserveBookWorkflowInput struct {
	BookID   int
	Reserver string
}

func ReserveBookWorkflow(ctx workflow.Context, input ReserveBookWorkflowInput) error {
	var book *models.Book
	err := workflow.ExecuteActivity(ctx, CheckBookAvailability, input.BookID).Get(ctx, &book)
	if err != nil {
		return err
	}


	if book.IsReserved {
		return fmt.Errorf("book is already reserved")
	}

	err = ReserveBook(ctx, input.BookID, input.Reserver)
	if err != nil {
		return err
	}

	return nil
}

func CheckBookAvailability(bookID int) (*models.Book, error) {
	book := &models.Book{
		ID:         uint(bookID),
		Title:      "Example Book",
		Author:     "John Doe",
		Quantity:   5,
		IsReserved: false, 
	}

	if book.ID == 0 {
		return nil, fmt.Errorf("book not found")
	}

	return book, nil
}

func ReserveBook(_ workflow.Context, bookID int, reserver string) error {
	fmt.Printf("Book with ID %d is reserved by %s\n", bookID, reserver)
	return nil
}
