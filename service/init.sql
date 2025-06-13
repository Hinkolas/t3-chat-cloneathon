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
        UNIQUE (chat_id, message_index)
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
        'Morning Routine Ideas',
        'claude-4-sonnet',
        0,
        0,
        1748791019
    ),
    (
        'test-chat-2',
        'user-123',
        'Travel Planning Tips',
        'qwen3',
        0,
        0,
        1748841019
    ),
    (
        'test-chat-3',
        'user-123',
        'Gardening for Beginners',
        'claude-4-sonnet',
        0,
        0,
        1748891019
    ),
    (
        'test-chat-4',
        'user-123',
        'Fitness at Home',
        'qwen3',
        1,
        0,
        1748941019
    ),
    (
        'test-chat-5',
        'user-123',
        'Budgeting Basics',
        'claude-4-sonnet',
        0,
        0,
        1748991019
    ),
    (
        'test-chat-6',
        'user-123',
        'Understanding Climate Change',
        'qwen3',
        0,
        0,
        1749041019
    ),
    (
        'test-chat-7',
        'user-123',
        'DIY Home Repairs',
        'claude-4-sonnet',
        0,
        0,
        1749091019
    ),
    (
        'test-chat-8',
        'user-123',
        'Healthy Sleep Habits',
        'qwen3',
        0,
        0,
        1749141019
    ),
    (
        'test-chat-9',
        'user-123',
        'Introduction to Investing',
        'claude-4-sonnet',
        0,
        0,
        1749191019
    ),
    (
        'test-chat-10',
        'user-123',
        'Cooking Quick Meals',
        'qwen3',
        0,
        0,
        1749241019
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
        'What are some effective morning routine ideas?',
        1748791019,
        1
    ),
    (
        'msg-1-2',
        'test-chat-1',
        'assistant',
        'Waking up early, hydrating, light exercise, and meditation can set a positive tone for your day.',
        1748791039,
        2
    ),
    (
        'msg-1-3',
        'test-chat-1',
        'user',
        'How long should a morning routine be?',
        1748791059,
        3
    ),
    (
        'msg-1-4',
        'test-chat-1',
        'assistant',
        'It varies for everyone, but even 15-30 minutes dedicated to self-care can make a difference.',
        1748791079,
        4
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
        'msg-2-1',
        'test-chat-2',
        'user',
        'What are essential travel planning tips?',
        1748841019,
        1
    ),
    (
        'msg-2-2',
        'test-chat-2',
        'assistant',
        'Start with a budget, research your destination, book accommodations and flights in advance, and pack light.',
        1748841039,
        2
    ),
    (
        'msg-2-3',
        'test-chat-2',
        'user',
        'How to find cheap flights?',
        1748841059,
        3
    ),
    (
        'msg-2-4',
        'test-chat-2',
        'assistant',
        'Use flight comparison websites, be flexible with your travel dates, consider flying mid-week, and look into budget airlines.',
        1748841079,
        4
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
        'msg-3-1',
        'test-chat-3',
        'user',
        'I want to start gardening. Where do I begin?',
        1748891019,
        1
    ),
    (
        'msg-3-2',
        'test-chat-3',
        'assistant',
        'Begin with easy-to-grow plants like herbs or leafy greens. Choose a spot with good sunlight and proper drainage.',
        1748891039,
        2
    ),
    (
        'msg-3-3',
        'test-chat-3',
        'user',
        'What tools do I need for gardening?',
        1748891059,
        3
    ),
    (
        'msg-3-4',
        'test-chat-3',
        'assistant',
        'Basic tools include gardening gloves, a hand trowel, a cultivator, and a watering can.',
        1748891079,
        4
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
        'msg-4-1',
        'test-chat-4',
        'user',
        'How can I stay fit at home without equipment?',
        1748941019,
        1
    ),
    (
        'msg-4-2',
        'test-chat-4',
        'assistant',
        'Bodyweight exercises like push-ups, squats, planks, and lunges are highly effective. You can also try jumping jacks and burpees for cardio.',
        1748941039,
        2
    ),
    (
        'msg-4-3',
        'test-chat-4',
        'user',
        'Are online workout videos helpful?',
        1748941059,
        3
    ),
    (
        'msg-4-4',
        'test-chat-4',
        'assistant',
        'Yes, many free online resources and apps offer guided home workouts suitable for all fitness levels. They can provide structure and motivation.',
        1748941079,
        4
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
        'msg-5-1',
        'test-chat-5',
        'user',
        'What are the basics of personal budgeting?',
        1748991019,
        1
    ),
    (
        'msg-5-2',
        'test-chat-5',
        'assistant',
        'Track your income and expenses, set financial goals, categorize your spending, and review your budget regularly to make adjustments.',
        1748991039,
        2
    ),
    (
        'msg-5-3',
        'test-chat-5',
        'user',
        'What is the 50/30/20 rule?',
        1748991059,
        3
    ),
    (
        'msg-5-4',
        'test-chat-5',
        'assistant',
        'It suggests allocating 50% of your income to needs, 30% to wants, and 20% to savings and debt repayment. It''s a simple guideline for budgeting.',
        1748991079,
        4
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
        'msg-6-1',
        'test-chat-6',
        'user',
        'Can you explain climate change in simple terms?',
        1749041019,
        1
    ),
    (
        'msg-6-2',
        'test-chat-6',
        'assistant',
        'Climate change refers to long-term shifts in temperatures and weather patterns, primarily caused by human activities leading to increased greenhouse gases in the atmosphere.',
        1749041039,
        2
    ),
    (
        'msg-6-3',
        'test-chat-6',
        'user',
        'What are the main impacts of climate change?',
        1749041059,
        3
    ),
    (
        'msg-6-4',
        'test-chat-6',
        'assistant',
        'Impacts include rising sea levels, more extreme weather events (heatwaves, floods, droughts), disruptions to ecosystems, and threats to food security.',
        1749041079,
        4
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
        'msg-7-1',
        'test-chat-7',
        'user',
        'What are some common DIY home repairs I can learn?',
        1749091019,
        1
    ),
    (
        'msg-7-2',
        'test-chat-7',
        'assistant',
        'Simple tasks like fixing a leaky faucet, patching small holes in walls, unclogging drains, and changing light fixtures are great starting points.',
        1749091039,
        2
    ),
    (
        'msg-7-3',
        'test-chat-7',
        'user',
        'Where can I find reliable DIY repair guides?',
        1749091059,
        3
    ),
    (
        'msg-7-4',
        'test-chat-7',
        'assistant',
        'YouTube tutorials, home improvement websites, and even your local library can provide step-by-step instructions and visual aids.',
        1749091079,
        4
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
        'msg-8-1',
        'test-chat-8',
        'user',
        'How can I develop healthier sleep habits?',
        1749141019,
        1
    ),
    (
        'msg-8-2',
        'test-chat-8',
        'assistant',
        'Establish a consistent sleep schedule, create a relaxing bedtime routine, ensure your bedroom is dark and cool, and avoid caffeine and screens before bed.',
        1749141039,
        2
    ),
    (
        'msg-8-3',
        'test-chat-8',
        'user',
        'Is napping good or bad?',
        1749141059,
        3
    ),
    (
        'msg-8-4',
        'test-chat-8',
        'assistant',
        'Short power naps (20-30 minutes) can be beneficial for alertness, but longer or irregular naps can disrupt nighttime sleep.',
        1749141079,
        4
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
        'msg-9-1',
        'test-chat-9',
        'user',
        'What should I know about investing for beginners?',
        1749191019,
        1
    ),
    (
        'msg-9-2',
        'test-chat-9',
        'assistant',
        'Start by understanding your financial goals and risk tolerance. Consider low-cost index funds or ETFs for diversification and long-term growth.',
        1749191039,
        2
    ),
    (
        'msg-9-3',
        'test-chat-9',
        'user',
        'Is it too risky to invest now?',
        1749191059,
        3
    ),
    (
        'msg-9-4',
        'test-chat-9',
        'assistant',
        'Investing always carries some risk, but time in the market generally outweighs timing the market. Start small and regularly contribute to mitigate risk over the long term.',
        1749191079,
        4
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
        'msg-10-1',
        'test-chat-10',
        'user',
        'What are some ideas for quick and healthy meals?',
        1749241019,
        1
    ),
    (
        'msg-10-2',
        'test-chat-10',
        'assistant',
        'Sheet pan dinners, stir-fries, loaded salads, and quick pasta dishes with plenty of vegetables are great options for busy weeknights.',
        1749241039,
        2
    ),
    (
        'msg-10-3',
        'test-chat-10',
        'user',
        'How can I meal prep efficiently?',
        1749241059,
        3
    ),
    (
        'msg-10-4',
        'test-chat-10',
        'assistant',
        'Dedicate a few hours on a weekend to chop vegetables, cook grains, and pre-cook proteins. Store components separately to mix and match during the week.',
        1749309419,
        4
    );

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
        'test-chat-11',
        'user-123',
        'Healthy Eating Tips',
        'claude-4-sonnet',
        0,
        0,
        1749309419
    ),
    (
        'test-chat-12',
        'user-123',
        'Learning a New Language',
        'qwen3',
        0,
        0,
        1749360000
    ),
    (
        'test-chat-13',
        'user-123',
        'Home Decor Ideas',
        'claude-4-sonnet',
        0,
        0,
        1749410000
    ),
    (
        'test-chat-14',
        'user-123',
        'Time Management',
        'qwen3',
        1,
        0,
        1749460000
    ),
    (
        'test-chat-15',
        'user-123',
        'Pet Care Advice',
        'claude-4-sonnet',
        0,
        0,
        1749510000
    ),
    (
        'test-chat-16',
        'user-123',
        'Understanding AI',
        'qwen3',
        0,
        0,
        1749560000
    ),
    (
        'test-chat-17',
        'user-123',
        'Coding Challenges',
        'claude-4-sonnet',
        0,
        0,
        1749610000
    ),
    (
        'test-chat-18',
        'user-123',
        'Mindfulness and Meditation',
        'qwen3',
        0,
        0,
        1749660000
    ),
    (
        'test-chat-19',
        'user-123',
        'Digital Photography',
        'claude-4-sonnet',
        0,
        0,
        1749710000
    ),
    (
        'test-chat-20',
        'user-123',
        'Car Maintenance Basics',
        'qwen3',
        0,
        0,
        1749760000
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
        'msg-11-1',
        'test-chat-11',
        'user',
        'What are some simple ways to eat healthier?',
        1749309419,
        1
    ),
    (
        'msg-11-2',
        'test-chat-11',
        'assistant',
        'Start by incorporating more fruits and vegetables into your daily meals. Aim for a variety of colors to ensure a wide range of nutrients.',
        1749309439,
        2
    ),
    (
        'msg-11-3',
        'test-chat-11',
        'user',
        'Should I cut out carbs completely?',
        1749309459,
        3
    ),
    (
        'msg-11-4',
        'test-chat-11',
        'assistant',
        'Not necessarily! Focus on complex carbohydrates like whole grains, oats, and brown rice, which provide sustained energy and fiber, rather than eliminating them entirely.',
        1749309479,
        4
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
        'msg-12-1',
        'test-chat-12',
        'user',
        'What''s the best way to learn a new language quickly?',
        1749360000,
        1
    ),
    (
        'msg-12-2',
        'test-chat-12',
        'assistant',
        'Immersion is key! Try to expose yourself to the language as much as possible through movies, music, and conversations with native speakers.',
        1749360020,
        2
    ),
    (
        'msg-12-3',
        'test-chat-12',
        'user',
        'Are language learning apps effective?',
        1749360040,
        3
    ),
    (
        'msg-12-4',
        'test-chat-12',
        'assistant',
        'Yes, apps like Duolingo or Babbel can be great for building vocabulary and basic grammar, but they are most effective when combined with other learning methods.',
        1749360060,
        4
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
        'msg-13-1',
        'test-chat-13',
        'user',
        'I want to redecorate my living room. Any budget-friendly ideas?',
        1749410000,
        1
    ),
    (
        'msg-13-2',
        'test-chat-13',
        'assistant',
        'Absolutely! Repainting walls, adding throw pillows and blankets, and rearranging existing furniture can make a big impact without spending much.',
        1749410020,
        2
    ),
    (
        'msg-13-3',
        'test-chat-13',
        'user',
        'What about lighting?',
        1749410040,
        3
    ),
    (
        'msg-13-4',
        'test-chat-13',
        'assistant',
        'Strategic lighting can transform a room. Consider adding floor lamps, table lamps, and even fairy lights to create different moods and highlight features.',
        1749410060,
        4
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
        'msg-14-1',
        'test-chat-14',
        'user',
        'I struggle with time management. How can I improve?',
        1749460000,
        1
    ),
    (
        'msg-14-2',
        'test-chat-14',
        'assistant',
        'Try the Pomodoro Technique: work for 25 minutes, then take a 5-minute break. This helps maintain focus and prevents burnout.',
        1749460020,
        2
    ),
    (
        'msg-14-3',
        'test-chat-14',
        'user',
        'How do I prioritize tasks?',
        1749460040,
        3
    ),
    (
        'msg-14-4',
        'test-chat-14',
        'assistant',
        'Use the Eisenhower Matrix: categorize tasks as urgent/important, important/not urgent, urgent/not important, and neither. Focus on the urgent and important ones first.',
        1749460060,
        4
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
        'msg-15-1',
        'test-chat-15',
        'user',
        'My cat is scratching the furniture. Any advice?',
        1749510000,
        1
    ),
    (
        'msg-15-2',
        'test-chat-15',
        'assistant',
        'Provide plenty of scratching posts of various textures and heights. Encourage use with catnip or toys, and trim your cat''s nails regularly.',
        1749510020,
        2
    ),
    (
        'msg-15-3',
        'test-chat-15',
        'user',
        'How often should I feed my dog?',
        1749510040,
        3
    ),
    (
        'msg-15-4',
        'test-chat-15',
        'assistant',
        'Most adult dogs do well with two meals a day, morning and evening. Puppies may need 3-4 smaller meals daily.',
        1749510060,
        4
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
        'msg-16-1',
        'test-chat-16',
        'user',
        'What exactly is Artificial Intelligence?',
        1749560000,
        1
    ),
    (
        'msg-16-2',
        'test-chat-16',
        'assistant',
        'Artificial Intelligence (AI) refers to the simulation of human intelligence in machines programmed to think like humans and mimic their actions, such as learning, problem-solving, and decision-making.',
        1749560020,
        2
    ),
    (
        'msg-16-3',
        'test-chat-16',
        'user',
        'What''s the difference between AI and Machine Learning?',
        1749560040,
        3
    ),
    (
        'msg-16-4',
        'test-chat-16',
        'assistant',
        'Machine Learning is a subset of AI that enables systems to learn from data without explicit programming. AI is the broader concept of creating intelligent machines, while ML is a method to achieve that intelligence.',
        1749560060,
        4
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
        'msg-17-1',
        'test-chat-17',
        'user',
        'Where can I find coding challenges for beginners?',
        1749610000,
        1
    ),
    (
        'msg-17-2',
        'test-chat-17',
        'assistant',
        'Websites like HackerRank, LeetCode (filter by easy), and freeCodeCamp offer a wide range of coding challenges perfect for beginners to practice their skills.',
        1749610020,
        2
    ),
    (
        'msg-17-3',
        'test-chat-17',
        'user',
        'What''s a good first project idea?',
        1749610040,
        3
    ),
    (
        'msg-17-4',
        'test-chat-17',
        'assistant',
        'A simple to-do list application, a basic calculator, or a rock-paper-scissors game are excellent first projects to solidify your foundational coding concepts.',
        1749610060,
        4
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
        'msg-18-1',
        'test-chat-18',
        'user',
        'How can I start practicing mindfulness?',
        1749660000,
        1
    ),
    (
        'msg-18-2',
        'test-chat-18',
        'assistant',
        'Begin with simple breathing exercises. Find a quiet spot, sit comfortably, and focus on your breath for 5-10 minutes. Notice the sensation of air entering and leaving your body.',
        1749660020,
        2
    ),
    (
        'msg-18-3',
        'test-chat-18',
        'user',
        'Is meditation difficult?',
        1749660040,
        3
    ),
    (
        'msg-18-4',
        'test-chat-18',
        'assistant',
        'It can feel challenging at first, but consistency is more important than perfection. Don''t worry if your mind wanders; gently bring your focus back to your breath. Guided meditations can be very helpful.',
        1749660060,
        4
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
        'msg-19-1',
        'test-chat-19',
        'user',
        'What are the basics of digital photography?',
        1749710000,
        1
    ),
    (
        'msg-19-2',
        'test-chat-19',
        'assistant',
        'The "exposure triangle" is fundamental: aperture (controls depth of field), shutter speed (controls motion blur), and ISO (controls sensitivity to light). Understanding how they interact is key.',
        1749710020,
        2
    ),
    (
        'msg-19-3',
        'test-chat-19',
        'user',
        'How important is composition?',
        1749710040,
        3
    ),
    (
        'msg-19-4',
        'test-chat-19',
        'assistant',
        'Composition is extremely important! It''s how you arrange elements within your frame. Rules like the Rule of Thirds, leading lines, and framing can dramatically improve your photos.',
        1749710060,
        4
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
        'msg-20-1',
        'test-chat-20',
        'user',
        'What basic car maintenance should I know?',
        1749760000,
        1
    ),
    (
        'msg-20-2',
        'test-chat-20',
        'assistant',
        'Regular oil changes, checking tire pressure, and ensuring your fluid levels (coolant, brake fluid, windshield washer fluid) are adequate are crucial for car longevity and safety.',
        1749760020,
        2
    ),
    (
        'msg-20-3',
        'test-chat-20',
        'user',
        'How often should I check my oil?',
        1749760040,
        3
    ),
    (
        'msg-20-4',
        'test-chat-20',
        'assistant',
        'It''s a good idea to check your engine oil level once a month, and definitely before a long road trip. This helps catch potential issues early.',
        1749827819,
        4
    );