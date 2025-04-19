CREATE DATABASE student_portal;

\c student_portal;

-- TODO: Create is_active field for admin, hod, fpc tablesso that the reviews table can be referenced

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE students (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    registration_number VARCHAR(20) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    official_mail_id VARCHAR(100) UNIQUE NOT NULL
);

CREATE TABLE student_submissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    registration_number VARCHAR(20) NOT NULL,
    name VARCHAR(100) NOT NULL,
    gender VARCHAR(10),
    cgpa DECIMAL(3, 2),
    backlogs INT,
    semester VARCHAR(10),
    official_mail_id VARCHAR(100) NOT NULL,
    mobile_number VARCHAR(15),
    department VARCHAR(50),
    section VARCHAR(10),
    offer_type VARCHAR(50),
    company_name VARCHAR(100),
    company_state VARCHAR(100), 
    company_city VARCHAR(100), 
    pincode VARCHAR(20), 
    hrd_email VARCHAR(100),
    hrd_number VARCHAR(15),
    offer_type_detail VARCHAR(100),
    package_ppo DECIMAL(10, 2),
    stipend_amount DECIMAL(10, 2),
    internship_start_date DATE NOT NULL,
    internship_end_date DATE NOT NULL,
    offer_letter_path VARCHAR(255),
    mail_copy_path VARCHAR(255),
    terms_accepted BOOLEAN DEFAULT FALSE,
    status VARCHAR(20) DEFAULT 'Pending',
    noc_type VARCHAR(20) NOT NULL,
    noc_path VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE fpc (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL, 
    email VARCHAR(100) NOT NULL UNIQUE,
    app_password TEXT,
    password_hash TEXT NOT NULL,
    department VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE fpc_reviews (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    submission_id UUID NOT NULL,
    fpc_id UUID,
    status VARCHAR(20) NOT NULL, 
    comments TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (submission_id) REFERENCES student_submissions(id),
    FOREIGN KEY (fpc_id) REFERENCES fpc(id) ON DELETE SET NULL
);

CREATE TABLE hod ( 
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL, 
    email VARCHAR(255) UNIQUE NOT NULL,
    app_password TEXT,
    password_hash TEXT NOT NULL,
    department VARCHAR(50) NOT NULL,
    role_type VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE hod_reviews (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    submission_id UUID NOT NULL,
    hod_id UUID,
    action VARCHAR(20) NOT NULL, 
    remarks TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (submission_id) REFERENCES student_submissions(id),
    FOREIGN KEY (hod_id) REFERENCES hod(id) ON DELETE SET NULL
);

CREATE TABLE admin (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL, 
    email VARCHAR(255) UNIQUE NOT NULL,
    app_password TEXT,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE office (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    department VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
