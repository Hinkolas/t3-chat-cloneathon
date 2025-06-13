CREATE TABLE
    chats (
        id TEXT PRIMARY KEY,
        user_id TEXT NOT NULL,
        title TEXT NOT NULL,
        model TEXT NOT NULL,
        pinned INTEGER NOT NULL,
        is_streaming INTEGER NOT NULL,
        created_at INTEGER NOT NULL
    );

CREATE TABLE
    messages (
        id TEXT PRIMARY KEY,
        chat_id TEXT NOT NULL REFERENCES chats (id) ON DELETE CASCADE,
        role TEXT NOT NULL,
        content TEXT NOT NULL,
        created_at INTEGER NOT NULL,
        message_index INTEGER NOT NULL,
        UNIQUE(chat_id, message_index)
    );

CREATE TABLE
    attachments (
        id TEXT PRIMARY KEY,
        message_id TEXT NOT NULL REFERENCES messages (id) ON DELETE CASCADE,
        name TEXT NOT NULL,
        type TEXT NOT NULL
    );

CREATE INDEX idx_messages_chat_id_created_at ON messages (chat_id, created_at);

INSERT INTO
    chats (
        id,
        user_id,
        title,
        model,
        pinned,
        is_streaming,
        created_at
    )
VALUES
    (
        'test-chat-1',
        'user-123',
        'Planning Weekend Trip',
        'qwen3',
        0,
        0,
        1703980800
    ),
    (
        'test-chat-2',
        'user-123',
        'Recipe Help',
        'claude-4-sonnet',
        0,
        0,
        1703981800
    ),
    (
        'test-chat-3',
        'user-123',
        'Career Advice',
        'claude-4-sonnet',
        1,
        0,
        1703982800
    );

INSERT INTO
    messages (
        id,
        chat_id,
        role,
        content,
        created_at,
        message_index
    )
VALUES
    (
        'msg-1-1',
        'test-chat-1',
        'user',
        'I need help planning a weekend trip to the mountains',
        1703980800,
        1
    ),
    (
        'msg-1-2',
        'test-chat-1',
        'assistant',
        'I''d be happy to help you plan your mountain getaway! What''s your budget and how many people are going?',
        1703980820,
        2
    ),
    (
        'msg-1-3',
        'test-chat-1',
        'user',
        'It''s just me and my partner, budget around $500 for the weekend',
        1703980840,
        3
    ),
    (
        'msg-1-4',
        'test-chat-1',
        'assistant',
        'Perfect! For $500, you could look into cabin rentals or a nice mountain lodge. Are you interested in hiking, skiing, or just relaxing?',
        1703980860,
        4
    ),
    (
        'msg-1-5',
        'test-chat-1',
        'user',
        'We love hiking and want to see some scenic views',
        1703980880,
        5
    ),
    (
        'msg-1-6',
        'test-chat-1',
        'assistant',
        'Great choice! I recommend checking out Blue Ridge Mountains or the Smokies. Both have excellent trail systems and breathtaking viewpoints.',
        1703980900,
        6
    ),
    (
        'msg-1-7',
        'test-chat-1',
        'user',
        'What about accommodation near Blue Ridge?',
        1703980920,
        7
    ),
    (
        'msg-1-8',
        'test-chat-1',
        'assistant',
        'There are several great options: mountain cabins on Airbnb ($150-200/night), or the Blue Ridge Lodge ($120/night). Both offer easy trail access.',
        1703980940,
        8
    ),
    (
        'msg-1-9',
        'test-chat-1',
        'user',
        'The lodge sounds perfect! Any must-see trails?',
        1703980960,
        9
    ),
    (
        'msg-1-10',
        'test-chat-1',
        'assistant',
        'Definitely try the Scenic Overlook Trail (moderate,3 miles) and Waterfall Loop (easy,2 miles). Both offer stunning photo opportunities!',
        1703980980,
        10
    ),
    (
        'msg-2-1',
        'test-chat-2',
        'user',
        'I want to make something special for dinner tonight but only have chicken and basic pantry items',
        1703981800,
        1
    ),
    (
        'msg-2-2',
        'test-chat-2',
        'assistant',
        'I can help you create something delicious! What pantry staples do you have? Rice, pasta, canned tomatoes, spices?',
        1703981820,
        2
    ),
    (
        'msg-2-3',
        'test-chat-2',
        'user',
        'I have rice, garlic, onions, olive oil, and some herbs like thyme and rosemary',
        1703981840,
        3
    ),
    (
        'msg-2-4',
        'test-chat-2',
        'assistant',
        'Perfect! You can make a wonderful one-pan chicken and rice dish. Season the chicken with thyme, rosemary, salt and pepper, then sear it first.',
        1703981860,
        4
    ),
    (
        'msg-2-5',
        'test-chat-2',
        'user',
        'How long should I sear the chicken?',
        1703981880,
        5
    ),
    (
        'msg-2-6',
        'test-chat-2',
        'assistant',
        'Sear for about3-4 minutes per side until golden brown. Remove chicken, then sauté diced onions and garlic in the same pan for2-3 minutes.',
        1703981900,
        6
    ),
    (
        'msg-2-7',
        'test-chat-2',
        'user',
        'Then I add the rice?',
        1703981920,
        7
    ),
    (
        'msg-2-8',
        'test-chat-2',
        'assistant',
        'Yes! Add1 cup rice, stir for1 minute, then add2 cups chicken broth or water. Place chicken back on top, cover and simmer18-20 minutes.',
        1703981940,
        8
    ),
    (
        'msg-2-9',
        'test-chat-2',
        'user',
        'What if I don''t have chicken broth?',
        1703981960,
        9
    ),
    (
        'msg-2-10',
        'test-chat-2',
        'assistant',
        'Water works fine! Just add extra salt and herbs. The chicken will release flavors as it cooks, making the rice delicious.',
        1703981980,
        10
    ),
    (
        'msg-3-6',
        'test-chat-3',
        'assistant',
        'Start with fundamentals: choose Python or JavaScript as your first language. Use free resources like freeCodeCamp, then build small projects to practice.',
        1703982900,
        1
    ),
    (
        'msg-3-7',
        'test-chat-3',
        'user',
        'Should I quit my job to focus on learning full-time?',
        1703982920,
        2
    ),
    (
        'msg-3-8',
        'test-chat-3',
        'assistant',
        'I''d recommend learning while working initially. Dedicate1-2 hours daily for6-12 months. Once you have solid basics and a portfolio, then consider bootcamps or transitioning.',
        1703982940,
        3
    ),
    (
        'msg-3-9',
        'test-chat-3',
        'user',
        'How long does it typically take to become job-ready?',
        1703982960,
        4
    ),
    (
        'msg-3-10',
        'test-chat-3',
        'assistant',
        'With consistent daily practice, most career changers become job-ready in12-18 months. Your timeline may be shorter given your professional experience and transferable skills.',
        1703982980,
        5
    );

INSERT INTO
    attachments (
        id,
        message_id,
        name,
        type
    )
VALUES
    (
        'att-1',
        'msg-1-10',
        'blue_ridge_trail_map.pdf',
        'application/pdf'
    ),
    (
        'att-2',
        'msg-1-10',
        'scenic_overlook_photo.jpg',
        'image/jpeg'
    ),
    (
        'att-3',
        'msg-2-8',
        'chicken_rice_recipe.pdf',
        'application/pdf'
    ),
    (
        'att-4',
        'msg-2-3',
        'pantry_ingredients.jpg',
        'image/jpeg'
    ),
    (
        'att-5',
        'msg-3-6',
        'programming_learning_roadmap.pdf',
        'application/pdf'
    ),
    (
        'att-6',
        'msg-3-8',
        'coding_bootcamp_comparison.xlsx',
        'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'
    ),
    (
        'att-7',
        'msg-1-8',
        'blue_ridge_lodge_brochure.pdf',
        'application/pdf'
    ),
    (
        'att-8',
        'msg-2-10',
        'finished_dish_photo.png',
        'image/png'
    );
