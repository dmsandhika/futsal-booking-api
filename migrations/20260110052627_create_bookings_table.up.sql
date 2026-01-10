CREATE TABLE bookings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    court_id UUID NOT NULL,
    booking_date DATE NOT NULL,
    time_slot VARCHAR(50) NOT NULL,
    customer_name VARCHAR(255) NOT NULL,
    customer_contact VARCHAR(100),
    customer_email VARCHAR(255),
    total_price INT NOT NULL,
    payment_type VARCHAR(50),
    amount_paid INT,
    remaining_amount INT,
    payment_status VARCHAR(50) DEFAULT 'unpaid',
    payment_deadline TIMESTAMP WITH TIME ZONE,
    status VARCHAR(50) DEFAULT 'pending',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (court_id) REFERENCES courts(id) ON DELETE CASCADE
);