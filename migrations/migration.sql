CREATE DATABASE student_portal;

\c student_portal;

CREATE TABLE student_submissions (
    id SERIAL PRIMARY KEY,
    registration_number VARCHAR(20) NOT NULL,
    name VARCHAR(100) NOT NULL,
    official_mail_id VARCHAR(100) NOT NULL,
    mobile_number VARCHAR(15),
    department VARCHAR(50),
    section VARCHAR(10),
    offer_type VARCHAR(50),
    company_name VARCHAR(100),
    company_state VARCHAR(100), 
    company_city VARCHAR(100), 
    pincode VARCHAR(20), 
    offer_type_detail VARCHAR(100),
    package_ppo DECIMAL(10, 2),
    stipend_amount DECIMAL(10, 2),
    internship_start_date DATE NOT NULL,
    internship_end_date DATE NOT NULL,
    offer_letter_path VARCHAR(255),
    mail_copy_path VARCHAR(255),
    terms_accepted BOOLEAN DEFAULT FALSE,
    status VARCHAR(20) DEFAULT 'Pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE reviewers (                       
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL, 
    email VARCHAR(100) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    department VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE reviewer_reviews (
    id SERIAL PRIMARY KEY,
    submission_id INT NOT NULL,
    reviewer_id INT NOT NULL,
    status VARCHAR(20) NOT NULL, -- Accepted, Rejected, Rework
    comments TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (submission_id) REFERENCES student_submissions(id),
    FOREIGN KEY (reviewer_id) REFERENCES reviewers(id)
);


CREATE TABLE hod ( 
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL, 
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    department VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE hod_reviews (
    id SERIAL PRIMARY KEY,
    submission_id INT NOT NULL,
    hod_id INT NOT NULL,
    action VARCHAR(20) NOT NULL, -- Approved, Rejected, Rework
    remarks TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (submission_id) REFERENCES student_submissions(id),
    FOREIGN KEY (hod_id) REFERENCES hod(id)
);

