-- ============================================
-- COMPREHENSIVE TEST DATA FOR PAYMENT SCHEDULE
-- WITH user_id and contact_id in debt_items
-- ============================================

DO $$
DECLARE
    v_user_id uuid;
    v_contact_id uuid;
    v_debt_id uuid;
    v_created_date timestamp;
    v_due_date timestamp;
    v_next_payment_date timestamp;
BEGIN
    -- Get valid user and contact
    SELECT uc.user_id, uc.contact_id 
    INTO v_user_id, v_contact_id
    FROM user_contacts uc
    INNER JOIN contacts c ON uc.contact_id = c.id
    LIMIT 1;
    
    IF v_user_id IS NULL OR v_contact_id IS NULL THEN
        RAISE EXCEPTION 'No valid user-contact relationship found. Please create a contact first.';
    END IF;

    -- ============================================
    -- Test Case 1: WEEKLY - Overdue with mixed payments
    -- ============================================
    v_debt_id := '11111111-1111-1111-1111-111111111111'::uuid;
    v_created_date := NOW() - INTERVAL '70 days';
    v_due_date := v_created_date + INTERVAL '70 days';
    v_next_payment_date := v_created_date + INTERVAL '28 days';
    
    INSERT INTO debt_lists (
        id, user_id, contact_id, debt_type, total_amount, installment_amount,
        total_payments_made, total_remaining_debt, currency, status,
        due_date, next_payment_date, installment_plan, number_of_payments,
        description, notes, created_at, updated_at
    ) VALUES (
        v_debt_id, v_user_id, v_contact_id, 'i_owe',
        500.00, 50.00, 150.00, 350.00, 'Php', 'overdue',
        v_due_date, v_next_payment_date,
        'weekly', 10,
        'Test WEEKLY installments - overdue',
        NULL, v_created_date, NOW()
    );
    
    -- Insert payments with user_id and contact_id
    INSERT INTO debt_items (id, debt_list_id, amount, currency, payment_date, payment_method, description, status, created_at, updated_at)
    SELECT 
        gen_random_uuid(),
        v_debt_id,
        payment_data.amount,
        'Php',
        payment_data.payment_date,
        payment_data.payment_method,
        payment_data.description,
        payment_data.status,
        payment_data.payment_date,
        payment_data.payment_date
    FROM (VALUES
        (50.00, v_created_date + INTERVAL '7 days', 'bank_transfer', 'Week 1 payment', 'completed'),
        (50.00, v_created_date + INTERVAL '14 days', 'cash', 'Week 2 payment', 'completed'),
        (50.00, v_created_date + INTERVAL '21 days', 'digital_wallet', 'Week 3 payment', 'completed'),
        (50.00, v_created_date + INTERVAL '28 days', 'bank_transfer', 'Week 4 payment - FAILED', 'failed')
    ) AS payment_data(amount, payment_date, payment_method, description, status);

    -- ============================================
    -- Test Case 2: MONTHLY - Active with some payments
    -- ============================================
    v_debt_id := '22222222-2222-2222-2222-222222222222'::uuid;
    v_created_date := NOW() - INTERVAL '3 months' - INTERVAL '28 days';
    v_due_date := v_created_date + INTERVAL '12 months';
    v_next_payment_date := v_created_date + INTERVAL '4 months';
    
    INSERT INTO debt_lists (
        id, user_id, contact_id, debt_type, total_amount, installment_amount,
        total_payments_made, total_remaining_debt, currency, status,
        due_date, next_payment_date, installment_plan, number_of_payments,
        description, notes, created_at, updated_at
    ) VALUES (
        v_debt_id, v_user_id, v_contact_id, 'owed_to_me',
        1200.00, 100.00, 300.00, 900.00, 'Php', 'active',
        v_due_date, v_next_payment_date,
        'monthly', 12,
        'Test MONTHLY installments - active',
        NULL, v_created_date, NOW()
    );
    
    INSERT INTO debt_items (id, debt_list_id, amount, currency, payment_date, payment_method, description, status, created_at, updated_at)
    SELECT 
        gen_random_uuid(),
        v_debt_id,
        payment_data.amount,
        'Php',
        payment_data.payment_date,
        payment_data.payment_method,
        payment_data.description,
        'completed',
        payment_data.payment_date,
        payment_data.payment_date
    FROM (VALUES
        (100.00, v_created_date + INTERVAL '1 month', 'bank_transfer', 'Month 1 payment'),
        (100.00, v_created_date + INTERVAL '2 months', 'bank_transfer', 'Month 2 payment'),
        (100.00, v_created_date + INTERVAL '3 months', 'cash', 'Month 3 payment')
    ) AS payment_data(amount, payment_date, payment_method, description);

    -- ============================================
    -- Test Case 3: BIWEEKLY - Overdue with pending payment
    -- ============================================
    v_debt_id := '33333333-3333-3333-3333-333333333333'::uuid;
    v_created_date := NOW() - INTERVAL '112 days';
    v_due_date := v_created_date + INTERVAL '112 days';
    v_next_payment_date := v_created_date + INTERVAL '42 days';
    
    INSERT INTO debt_lists (
        id, user_id, contact_id, debt_type, total_amount, installment_amount,
        total_payments_made, total_remaining_debt, currency, status,
        due_date, next_payment_date, installment_plan, number_of_payments,
        description, notes, created_at, updated_at
    ) VALUES (
        v_debt_id, v_user_id, v_contact_id, 'i_owe',
        800.00, 100.00, 200.00, 600.00, 'Php', 'overdue',
        v_due_date, v_next_payment_date,
        'biweekly', 8,
        'Test BIWEEKLY installments - overdue with pending',
        NULL, v_created_date, NOW()
    );
    
    INSERT INTO debt_items (id, debt_list_id, amount, currency, payment_date, payment_method, description, status, created_at, updated_at)
    SELECT 
        gen_random_uuid(),
        v_debt_id,
        payment_data.amount,
        'Php',
        payment_data.payment_date,
        payment_data.payment_method,
        payment_data.description,
        payment_data.status,
        payment_data.payment_date,
        payment_data.payment_date
    FROM (VALUES
        (100.00, v_created_date + INTERVAL '14 days', 'bank_transfer', 'Biweekly 1', 'completed'),
        (100.00, v_created_date + INTERVAL '28 days', 'cash', 'Biweekly 2', 'completed'),
        (100.00, v_created_date + INTERVAL '42 days', 'digital_wallet', 'Biweekly 3 - PENDING', 'pending')
    ) AS payment_data(amount, payment_date, payment_method, description, status);

    -- ============================================
    -- Test Case 4: QUARTERLY - Active, early payments
    -- ============================================
    v_debt_id := '44444444-4444-4444-4444-444444444444'::uuid;
    v_created_date := NOW() - INTERVAL '3 months';
    v_due_date := v_created_date + INTERVAL '12 months';
    v_next_payment_date := v_created_date + INTERVAL '6 months';
    
    INSERT INTO debt_lists (
        id, user_id, contact_id, debt_type, total_amount, installment_amount,
        total_payments_made, total_remaining_debt, currency, status,
        due_date, next_payment_date, installment_plan, number_of_payments,
        description, notes, created_at, updated_at
    ) VALUES (
        v_debt_id, v_user_id, v_contact_id, 'owed_to_me',
        2000.00, 500.00, 500.00, 1500.00, 'Php', 'active',
        v_due_date, v_next_payment_date,
        'quarterly', 4,
        'Test QUARTERLY installments - active',
        NULL, v_created_date, NOW()
    );
    
    INSERT INTO debt_items (id, debt_list_id, amount, currency, payment_date, payment_method, description, status, created_at, updated_at)
    VALUES (
        gen_random_uuid(),
        v_debt_id,
        500.00,
        'Php',
        v_created_date + INTERVAL '3 months',
        'bank_transfer',
        'Quarter 1 payment',
        'completed',
        v_created_date + INTERVAL '3 months',
        v_created_date + INTERVAL '3 months'
    );

    -- ============================================
    -- Test Case 5: ONETIME - Overdue, no payments
    -- ============================================
    v_debt_id := '55555555-5555-5555-5555-555555555555'::uuid;
    v_created_date := NOW() - INTERVAL '60 days';
    v_due_date := NOW() - INTERVAL '30 days';
    v_next_payment_date := v_due_date;
    
    INSERT INTO debt_lists (
        id, user_id, contact_id, debt_type, total_amount, installment_amount,
        total_payments_made, total_remaining_debt, currency, status,
        due_date, next_payment_date, installment_plan, number_of_payments,
        description, notes, created_at, updated_at
    ) VALUES (
        v_debt_id, v_user_id, v_contact_id, 'i_owe',
        5000.00, 5000.00, 0.00, 5000.00, 'Php', 'overdue',
        v_due_date, v_next_payment_date,
        'onetime', 1,
        'Test ONETIME payment - overdue, no payments',
        NULL, v_created_date, NOW()
    );
    -- No payments for this one

    -- ============================================
    -- Test Case 6: YEARLY - Active, long term
    -- ============================================
    v_debt_id := '66666666-6666-6666-6666-666666666666'::uuid;
    v_created_date := NOW() - INTERVAL '2 years' - INTERVAL '1 month';
    v_due_date := v_created_date + INTERVAL '5 years';
    v_next_payment_date := v_created_date + INTERVAL '3 years';
    
    INSERT INTO debt_lists (
        id, user_id, contact_id, debt_type, total_amount, installment_amount,
        total_payments_made, total_remaining_debt, currency, status,
        due_date, next_payment_date, installment_plan, number_of_payments,
        description, notes, created_at, updated_at
    ) VALUES (
        v_debt_id, v_user_id, v_contact_id, 'i_owe',
        10000.00, 2000.00, 4000.00, 6000.00, 'Php', 'active',
        v_due_date, v_next_payment_date,
        'yearly', 5,
        'Test YEARLY installments - active',
        NULL, v_created_date, NOW()
    );
    
    INSERT INTO debt_items (id, debt_list_id, amount, currency, payment_date, payment_method, description, status, created_at, updated_at)
    SELECT 
        gen_random_uuid(),
        v_debt_id,
        payment_data.amount,
        'Php',
        payment_data.payment_date,
        'bank_transfer',
        payment_data.description,
        'completed',
        payment_data.payment_date,
        payment_data.payment_date
    FROM (VALUES
        (2000.00, v_created_date + INTERVAL '1 year', 'Year 1 payment'),
        (2000.00, v_created_date + INTERVAL '2 years', 'Year 2 payment')
    ) AS payment_data(amount, payment_date, description);

    -- ============================================
    -- Test Case 7: MONTHLY - Partially paid, overdue
    -- ============================================
    v_debt_id := '77777777-7777-7777-7777-777777777777'::uuid;
    v_created_date := NOW() - INTERVAL '5 months';
    v_due_date := NOW();
    v_next_payment_date := v_created_date + INTERVAL '3 months';
    
    INSERT INTO debt_lists (
        id, user_id, contact_id, debt_type, total_amount, installment_amount,
        total_payments_made, total_remaining_debt, currency, status,
        due_date, next_payment_date, installment_plan, number_of_payments,
        description, notes, created_at, updated_at
    ) VALUES (
        v_debt_id, v_user_id, v_contact_id, 'owed_to_me',
        1000.00, 200.00, 550.00, 450.00, 'Php', 'overdue',
        v_due_date, v_next_payment_date,
        'monthly', 5,
        'Test MONTHLY - partially paid installment',
        NULL, v_created_date, NOW()
    );
    
    INSERT INTO debt_items (id, debt_list_id, amount, currency, payment_date, payment_method, description, status, created_at, updated_at)
    SELECT 
        gen_random_uuid(),
        v_debt_id,
        payment_data.amount,
        'Php',
        payment_data.payment_date,
        payment_data.payment_method,
        payment_data.description,
        'completed',
        payment_data.payment_date,
        payment_data.payment_date
    FROM (VALUES
        (200.00, v_created_date + INTERVAL '1 month', 'bank_transfer', 'Month 1 payment'),
        (200.00, v_created_date + INTERVAL '2 months', 'cash', 'Month 2 payment'),
        (150.00, v_created_date + INTERVAL '3 months', 'digital_wallet', 'Month 3 partial payment')
    ) AS payment_data(amount, payment_date, payment_method, description);

    -- ============================================
    -- Test Case 8: WEEKLY - Near completion
    -- ============================================
    v_debt_id := '88888888-8888-8888-8888-888888888888'::uuid;
    v_created_date := NOW() - INTERVAL '56 days';
    v_due_date := NOW();
    v_next_payment_date := v_created_date + INTERVAL '56 days';
    
    INSERT INTO debt_lists (
        id, user_id, contact_id, debt_type, total_amount, installment_amount,
        total_payments_made, total_remaining_debt, currency, status,
        due_date, next_payment_date, installment_plan, number_of_payments,
        description, notes, created_at, updated_at
    ) VALUES (
        v_debt_id, v_user_id, v_contact_id, 'i_owe',
        400.00, 50.00, 350.00, 50.00, 'Php', 'active',
        v_due_date, v_next_payment_date,
        'weekly', 8,
        'Test WEEKLY - almost complete',
        NULL, v_created_date, NOW()
    );
    
    INSERT INTO debt_items (id, debt_list_id, amount, currency, payment_date, payment_method, status, created_at, updated_at)
    SELECT 
        gen_random_uuid(),
        v_debt_id,
        50.00,
        'Php',
        v_created_date + (INTERVAL '1 day' * n * 7),
        CASE (n % 3)
            WHEN 0 THEN 'bank_transfer'
            WHEN 1 THEN 'cash'
            ELSE 'digital_wallet'
        END,
        'completed',
        v_created_date + (INTERVAL '1 day' * n * 7),
        v_created_date + (INTERVAL '1 day' * n * 7)
    FROM generate_series(1, 7) AS n;

    -- ============================================
    -- Test Case 9: SETTLED debt (fully paid)
    -- ============================================
    v_debt_id := '99999999-9999-9999-9999-999999999999'::uuid;
    v_created_date := NOW() - INTERVAL '60 days';
    v_due_date := NOW() - INTERVAL '15 days';
    v_next_payment_date := v_due_date;
    
    INSERT INTO debt_lists (
        id, user_id, contact_id, debt_type, total_amount, installment_amount,
        total_payments_made, total_remaining_debt, currency, status,
        due_date, next_payment_date, installment_plan, number_of_payments,
        description, notes, created_at, updated_at
    ) VALUES (
        v_debt_id, v_user_id, v_contact_id, 'i_owe',
        300.00, 100.00, 300.00, 0.00, 'Php', 'settled',
        v_due_date, v_next_payment_date,
        'monthly', 3,
        'Test MONTHLY - SETTLED (fully paid)',
        NULL, v_created_date, NOW()
    );
    
    INSERT INTO debt_items (id, debt_list_id, amount, currency, payment_date, payment_method, description, status, created_at, updated_at)
    SELECT 
        gen_random_uuid(),
        v_debt_id,
        payment_data.amount,
        'Php',
        payment_data.payment_date,
        payment_data.payment_method,
        payment_data.description,
        'completed',
        payment_data.payment_date,
        payment_data.payment_date
    FROM (VALUES
        (100.00, v_created_date + INTERVAL '1 month', 'bank_transfer', 'Final payment 1/3'),
        (100.00, v_created_date + INTERVAL '2 months', 'cash', 'Final payment 2/3'),
        (100.00, v_created_date + INTERVAL '45 days', 'digital_wallet', 'Final payment 3/3')
    ) AS payment_data(amount, payment_date, payment_method, description);

    -- ============================================
    -- Test Case 10: ARCHIVED debt
    -- ============================================
    v_debt_id := 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa'::uuid;
    v_created_date := NOW() - INTERVAL '1 year';
    v_due_date := NOW() - INTERVAL '6 months';
    v_next_payment_date := v_due_date;
    
    INSERT INTO debt_lists (
        id, user_id, contact_id, debt_type, total_amount, installment_amount,
        total_payments_made, total_remaining_debt, currency, status,
        due_date, next_payment_date, installment_plan, number_of_payments,
        description, notes, created_at, updated_at
    ) VALUES (
        v_debt_id, v_user_id, v_contact_id, 'owed_to_me',
        1000.00, 1000.00, 0.00, 1000.00, 'Php', 'archived',
        v_due_date, v_next_payment_date,
        'onetime', 1,
        'Test ONETIME - ARCHIVED (old debt)',
        'User archived this debt - no longer pursuing payment', v_created_date, NOW()
    );

    RAISE NOTICE '========================================';
    RAISE NOTICE '✓ Test data created successfully!';
    RAISE NOTICE '========================================';
    RAISE NOTICE '10 test scenarios created';
    RAISE NOTICE 'User ID: %', v_user_id;
    RAISE NOTICE 'Contact ID: %', v_contact_id;
    RAISE NOTICE '========================================';
    RAISE NOTICE 'All debt_items include references to debt_list';
    RAISE NOTICE 'Use JOIN to get user_id and contact_id:';
    RAISE NOTICE '  SELECT di.*, dl.user_id, dl.contact_id';
    RAISE NOTICE '  FROM debt_items di';
    RAISE NOTICE '  JOIN debt_lists dl ON di.debt_list_id = dl.id';
    RAISE NOTICE '========================================';
END $$;


-- ============================================
-- DELETE BY SPECIFIC UUIDs
-- ============================================

DO $$
DECLARE
    v_test_debt_ids uuid[] := ARRAY[
        '11111111-1111-1111-1111-111111111111'::uuid,
        '22222222-2222-2222-2222-222222222222'::uuid,
        '33333333-3333-3333-3333-333333333333'::uuid,
        '44444444-4444-4444-4444-444444444444'::uuid,
        '55555555-5555-5555-5555-555555555555'::uuid,
        '66666666-6666-6666-6666-666666666666'::uuid,
        '77777777-7777-7777-7777-777777777777'::uuid,
        '88888888-8888-8888-8888-888888888888'::uuid,
		'99999999-9999-9999-9999-999999999999'::uuid,
		'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa'::uuid
    ];
    v_debt_items_deleted int;
    v_debt_lists_deleted int;
BEGIN
    -- Delete debt_items for these debts
    DELETE FROM debt_items 
    WHERE debt_list_id = ANY(v_test_debt_ids);
    GET DIAGNOSTICS v_debt_items_deleted = ROW_COUNT;
    
    -- Delete the debt_lists
    DELETE FROM debt_lists 
    WHERE id = ANY(v_test_debt_ids);
    GET DIAGNOSTICS v_debt_lists_deleted = ROW_COUNT;
    
    RAISE NOTICE '✓ Deleted % debt items', v_debt_items_deleted;
    RAISE NOTICE '✓ Deleted % debt lists', v_debt_lists_deleted;
END $$;