CREATE TABLE
    chats (
        id TEXT PRIMARY KEY,
        user_id TEXT NOT NULL,
        title TEXT NOT NULL,
        model TEXT NOT NULL,
        is_pinned INTEGER NOT NULL,
        is_streaming INTEGER NOT NULL,
        last_message_at INTEGER NOT NULL,
        created_at INTEGER NOT NULL,
        updated_at INTEGER NOT NULL
    );

CREATE TABLE
    messages (
        id TEXT PRIMARY KEY,
        chat_id TEXT NOT NULL REFERENCES chats (id) ON DELETE CASCADE,
        role TEXT NOT NULL,
        model TEXT NOT NULL,
        content TEXT NOT NULL,
        reasoning TEXT NOT NULL,
        created_at INTEGER NOT NULL,
        updated_at INTEGER NOT NULL
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
        is_pinned,
        is_streaming,
        last_message_at,
        created_at,
        updated_at
    )
VALUES
    (
        'test-chat-1',
        'user-123',
        'Morning Routine Ideas',
        'claude-4-sonnet',
        0,
        0,
        1748791079,
        1748791019,
        1748791079
    ),
    (
        'test-chat-2',
        'user-123',
        'Travel Planning Tips',
        'qwen3',
        0,
        0,
        1748841079,
        1748841019,
        1748841079
    ),
    (
        'test-chat-3',
        'user-123',
        'Gardening for Beginners',
        'claude-4-sonnet',
        0,
        0,
        1748891079,
        1748891019,
        1748891079
    ),
    (
        'test-chat-4',
        'user-123',
        'Fitness at Home',
        'qwen3',
        1,
        0,
        1748941079,
        1748941019,
        1748941079
    ),
    (
        'test-chat-5',
        'user-123',
        'Budgeting Basics',
        'claude-4-sonnet',
        0,
        0,
        1748991079,
        1748991019,
        1748991079
    ),
    (
        'test-chat-6',
        'user-123',
        'Understanding Climate Change',
        'qwen3',
        0,
        0,
        1749041079,
        1749041019,
        1749041079
    ),
    (
        'test-chat-7',
        'user-123',
        'DIY Home Repairs',
        'claude-4-sonnet',
        0,
        0,
        1749091079,
        1749091019,
        1749091079
    ),
    (
        'test-chat-8',
        'user-123',
        'Healthy Sleep Habits',
        'qwen3',
        0,
        0,
        1749141079,
        1749141019,
        1749141079
    ),
    (
        'test-chat-9',
        'user-123',
        'Introduction to Investing',
        'claude-4-sonnet',
        0,
        0,
        1749191079,
        1749191019,
        1749191079
    ),
    (
        'test-chat-10',
        'user-123',
        'Cooking Quick Meals',
        'qwen3',
        0,
        0,
        1749241019,
        1749241019,
        1749241019
    );

INSERT INTO
    messages (
        id,
        chat_id,
        role,
        model,
        content,
        reasoning,
        created_at,
        updated_at
    )
VALUES
    (
        'msg-1-1',
        'test-chat-1',
        'user',
        'claude-4-sonnet',
        'What are some effective morning routine ideas?',
        '',
        1748791019,
        1748791019
    ),
    (
        'msg-1-2',
        'test-chat-1',
        'assistant',
        'claude-4-sonnet',
        'Waking up early, hydrating, light exercise, and meditation can set a positive tone for your day.',
        '',
        1748791039,
        1748791039
    ),
    (
        'msg-1-3',
        'test-chat-1',
        'user',
        'claude-4-sonnet',
        'How long should a morning routine be?',
        '',
        1748791059,
        1748791059
    ),
    (
        'msg-1-4',
        'test-chat-1',
        'assistant',
        'claude-4-sonnet',
        'It varies for everyone, but even 15-30 minutes dedicated to self-care can make a difference.',
        '',
        1748791079,
        1748791079
    );

INSERT INTO
    messages (
        id,
        chat_id,
        role,
        model,
        content,
        reasoning,
        created_at,
        updated_at
    )
VALUES
    (
        'msg-2-1',
        'test-chat-2',
        'user',
        'qwen3',
        'What are essential travel planning tips?',
        '',
        1748841019,
        1748841019
    ),
    (
        'msg-2-2',
        'test-chat-2',
        'assistant',
        'qwen3',
        'Start with a budget, research your destination, book accommodations and flights in advance, and pack light.',
        'The user asked for essential travel planning tips. I will focus on key, actionable steps that are generally applicable to most travel scenarios. This includes financial preparation, destination knowledge, logistical arrangements (flights/accommodation), and practical packing advice to cover the main phases of planning a trip.',
        1748841039,
        1748841039
    ),
    (
        'msg-2-3',
        'test-chat-2',
        'user',
        'qwen3',
        'How to find cheap flights?',
        '',
        1748841059,
        1748841059
    ),
    (
        'msg-2-4',
        'test-chat-2',
        'assistant',
        'qwen3',
        'Use flight comparison websites, be flexible with your travel dates, consider flying mid-week, and look into budget airlines.',
        'The user asked for ways to find cheap flights. I will provide a range of common and effective strategies. These include using tools designed for price comparison, being open to less popular travel times (flexibility, mid-week), and considering cost-effective airline types (budget airlines).',
        1748841079,
        1748841079
    );

INSERT INTO
    messages (
        id,
        chat_id,
        role,
        model,
        content,
        reasoning,
        created_at,
        updated_at
    )
VALUES
    (
        'msg-3-1',
        'test-chat-3',
        'user',
        'claude-4-sonnet',
        'I want to start gardening. Where do I begin?',
        '',
        1748891019,
        1748891019
    ),
    (
        'msg-3-2',
        'test-chat-3',
        'assistant',
        'claude-4-sonnet',
        'Begin with easy-to-grow plants like herbs or leafy greens. Choose a spot with good sunlight and proper drainage.',
        '',
        1748891039,
        1748891039
    ),
    (
        'msg-3-3',
        'test-chat-3',
        'user',
        'claude-4-sonnet',
        'What tools do I need for gardening?',
        '',
        1748891059,
        1748891059
    ),
    (
        'msg-3-4',
        'test-chat-3',
        'assistant',
        'claude-4-sonnet',
        'Basic tools include gardening gloves, a hand trowel, a cultivator, and a watering can.',
        '',
        1748891079,
        1748891079
    );

INSERT INTO
    messages (
        id,
        chat_id,
        role,
        model,
        content,
        reasoning,
        created_at,
        updated_at
    )
VALUES
    (
        'msg-4-1',
        'test-chat-4',
        'user',
        'qwen3',
        'How can I stay fit at home without equipment?',
        '',
        1748941019,
        1748941019
    ),
    (
        'msg-4-2',
        'test-chat-4',
        '_assistant',
        'qwen3',
        'Bodyweight exercises like push-ups, squats, planks, and lunges are highly effective. You can also try jumping jacks and burpees for cardio.',
        'The user is asking for fitness methods without equipment. I will focus on bodyweight exercises as they are accessible and effective for both strength and cardiovascular fitness, requiring no external tools.',
        1748941039,
        1748941039
    ),
    (
        'msg-4-3',
        'test-chat-4',
        'user',
        'qwen3',
        'Are online workout videos helpful?',
        '',
        1748941059,
        1748941059
    ),
    (
        'msg-4-4',
        'test-chat-4',
        'assistant',
        'qwen3',
        'Yes, many free online resources and apps offer guided home workouts suitable for all fitness levels. They can provide structure and motivation.',
        'The user is asking about the utility of online workout videos. I will confirm their helpfulness and explain why, focusing on aspects like accessibility, guidance, structure, and motivational benefits, which are key advantages of such resources.',
        1748941079,
        1748941079
    );

INSERT INTO
    messages (
        id,
        chat_id,
        role,
        model,
        content,
        reasoning,
        created_at,
        updated_at
    )
VALUES
    (
        'msg-5-1',
        'test-chat-5',
        'user',
        'claude-4-sonnet',
        'What are the basics of personal budgeting?',
        '',
        1748991019,
        1748991019
    ),
    (
        'msg-5-2',
        'test-chat-5',
        'assistant',
        'claude-4-sonnet',
        'Track your income and expenses, set financial goals, categorize your spending, and review your budget regularly to make adjustments.',
        '',
        1748991039,
        1748991039
    ),
    (
        'msg-5-3',
        'test-chat-5',
        'user',
        'claude-4-sonnet',
        'What is the 50/30/20 rule?',
        '',
        1748991059,
        1748991059
    ),
    (
        'msg-5-4',
        'test-chat-5',
        'assistant',
        'claude-4-sonnet',
        'It suggests allocating 50% of your income to needs, 30% to wants, and 20% to savings and debt repayment. It''s a simple guideline for budgeting.',
        '',
        1748991079,
        1748991079
    );

INSERT INTO
    messages (
        id,
        chat_id,
        role,
        model,
        content,
        reasoning,
        created_at,
        updated_at
    )
VALUES
    (
        'msg-6-1',
        'test-chat-6',
        'user',
        'qwen3',
        'Can you explain climate change in simple terms?',
        '',
        1749041019,
        1749041019
    ),
    (
        'msg-6-2',
        'test-chat-6',
        'assistant',
        'qwen3',
        'Climate change refers to long-term shifts in temperatures and weather patterns, primarily caused by human activities leading to increased greenhouse gases in the atmosphere.',
        'The user is asking for a simple explanation of climate change. I will define it broadly by its core characteristics (long-term shifts in temperature/weather) and identify its main driver (human activities and resulting greenhouse gases) to provide a concise and understandable overview.',
        1749041039,
        1749041039
    ),
    (
        'msg-6-3',
        'test-chat-6',
        'user',
        'qwen3',
        'What are the main impacts of climate change?',
        '',
        1749041059,
        1749041059
    ),
    (
        'msg-6-4',
        'test-chat-6',
        'assistant',
        'qwen3',
        'Impacts include rising sea levels, more extreme weather events (heatwaves, floods, droughts), disruptions to ecosystems, and threats to food security.',
        'The user is asking for the main impacts of climate change. I will list several diverse and significant consequences, covering environmental effects (sea levels, extreme weather, ecosystems) and direct societal impacts (food security), to provide a comprehensive yet concise answer.',
        1749041079,
        1749041079
    );

INSERT INTO
    messages (
        id,
        chat_id,
        role,
        model,
        content,
        reasoning,
        created_at,
        updated_at
    )
VALUES
    (
        'msg-7-1',
        'test-chat-7',
        'user',
        'claude-4-sonnet',
        'What are some common DIY home repairs I can learn?',
        '',
        1749091019,
        1749091019
    ),
    (
        'msg-7-2',
        'test-chat-7',
        'assistant',
        'claude-4-sonnet',
        'Simple tasks like fixing a leaky faucet, patching small holes in walls, unclogging drains, and changing light fixtures are great starting points.',
        '',
        1749091039,
        1749091039
    ),
    (
        'msg-7-3',
        'test-chat-7',
        'user',
        'claude-4-sonnet',
        'Where can I find reliable DIY repair guides?',
        '',
        1749091059,
        1749091059
    ),
    (
        'msg-7-4',
        'test-chat-7',
        'assistant',
        'claude-4-sonnet',
        'YouTube tutorials, home improvement websites, and even your local library can provide step-by-step instructions and visual aids.',
        '',
        1749091079,
        1749091079
    );

INSERT INTO
    messages (
        id,
        chat_id,
        role,
        model,
        content,
        reasoning,
        created_at,
        updated_at
    )
VALUES
    (
        'msg-8-1',
        'test-chat-8',
        'user',
        'qwen3',
        'How can I develop healthier sleep habits?',
        '',
        1749141019,
        1749141019
    ),
    (
        'msg-8-2',
        'test-chat-8',
        'assistant',
        'qwen3',
        'Establish a consistent sleep schedule, create a relaxing bedtime routine, ensure your bedroom is dark and cool, and avoid caffeine and screens before bed.',
        'The user is asking for ways to develop healthier sleep habits. I will suggest a comprehensive set of practices that address different aspects of sleep hygiene: consistency in timing, pre-sleep relaxation, optimizing the sleep environment, and avoiding stimulants/disruptors.',
        1749141039,
        1749141039
    ),
    (
        'msg-8-3',
        'test-chat-8',
        'user',
        'qwen3',
        'Is napping good or bad?',
        '',
        1749141059,
        1749141059
    ),
    (
        'msg-8-4',
        'test-chat-8',
        'assistant',
        'qwen3',
        'Short power naps (20-30 minutes) can be beneficial for alertness, but longer or irregular naps can disrupt nighttime sleep.',
        'The user is asking about the benefits/downsides of napping. I will provide a nuanced answer, differentiating between beneficial "power naps" (short duration) and potentially disruptive "longer or irregular naps," explaining the impact on alertness versus nighttime sleep.',
        1749141079,
        1749141079
    );

INSERT INTO
    messages (
        id,
        chat_id,
        role,
        model,
        content,
        reasoning,
        created_at,
        updated_at
    )
VALUES
    (
        'msg-9-1',
        'test-chat-9',
        'user',
        'claude-4-sonnet',
        'What should I know about investing for beginners?',
        '',
        1749191019,
        1749191019
    ),
    (
        'msg-9-2',
        'test-chat-9',
        'assistant',
        'claude-4-sonnet',
        'Start by understanding your financial goals and risk tolerance. Consider low-cost index funds or ETFs for diversification and long-term growth.',
        '',
        1749191039,
        1749191039
    ),
    (
        'msg-9-3',
        'test-chat-9',
        'user',
        'claude-4-sonnet',
        'Is it too risky to invest now?',
        '',
        1749191059,
        1749191059
    ),
    (
        'msg-9-4',
        'test-chat-9',
        'assistant',
        'claude-4-sonnet',
        'Investing always carries some risk, but time in the market generally outweighs timing the market. Start small and regularly contribute to mitigate risk over the long term.',
        '',
        1749191079,
        1749191079
    );

INSERT INTO
    messages (
        id,
        chat_id,
        role,
        model,
        content,
        reasoning,
        created_at,
        updated_at
    )
VALUES
    (
        'msg-10-1',
        'test-chat-10',
        'user',
        'qwen3',
        'What are some ideas for quick and healthy meals?',
        '',
        1749241019,
        1749241019
    ),
    (
        'msg-10-2',
        'test-chat-10',
        'assistant',
        'qwen3',
        'Sheet pan dinners, stir-fries, loaded salads, and quick pasta dishes with plenty of vegetables are great options for busy weeknights.',
        'The user is asking for quick and healthy meal ideas. I will suggest categories of meals known for being both fast to prepare and nutritious, focusing on methods that simplify cooking and incorporate plenty of vegetables, suitable for busy schedules.',
        1749241039,
        1749241039
    ),
    (
        'msg-10-3',
        'test-chat-10',
        'user',
        'qwen3',
        'How can I meal prep efficiently?',
        '',
        1749241059,
        1749241059
    ),
    (
        'msg-10-4',
        'test-chat-10',
        'assistant',
        'qwen3',
        'Dedicate a few hours on a weekend to chop vegetables, cook grains, and pre-cook proteins. Store components separately to mix and match during the week.',
        'The user is asking for efficient meal prep strategies. I will outline a practical approach that involves dedicated time for batch preparation of core components (vegetables, grains, proteins) and emphasizes separate storage to maintain freshness and allow for versatile meal assembly throughout the week.',
        1749309419,
        1749309419
    );

INSERT INTO
    chats (
        id,
        user_id,
        title,
        model,
        is_pinned,
        is_streaming,
        last_message_at,
        created_at,
        updated_at
    )
VALUES
    (
        'test-chat-11',
        'user-123',
        'Healthy Eating Tips',
        'claude-4-sonnet',
        0,
        0,
        1749309479,
        1749309419,
        1749309479
    ),
    (
        'test-chat-12',
        'user-123',
        'Learning a New Language',
        'qwen3',
        0,
        0,
        1749360060,
        1749360000,
        1749360060
    ),
    (
        'test-chat-13',
        'user-123',
        'Home Decor Ideas',
        'claude-4-sonnet',
        0,
        0,
        1749410060,
        1749410000,
        1749410060
    ),
    (
        'test-chat-14',
        'user-123',
        'Time Management',
        'qwen3',
        1,
        0,
        1749460060,
        1749460000,
        1749460060
    ),
    (
        'test-chat-15',
        'user-123',
        'Pet Care Advice',
        'claude-4-sonnet',
        0,
        0,
        1749510060,
        1749510000,
        1749510060
    ),
    (
        'test-chat-16',
        'user-123',
        'Understanding AI',
        'qwen3',
        0,
        0,
        1749560060,
        1749560000,
        1749560060
    ),
    (
        'test-chat-17',
        'user-123',
        'Coding Challenges',
        'claude-4-sonnet',
        0,
        0,
        1749610060,
        1749610000,
        1749610060
    ),
    (
        'test-chat-18',
        'user-123',
        'Mindfulness and Meditation',
        'qwen3',
        0,
        0,
        1749660060,
        1749660000,
        1749660060
    ),
    (
        'test-chat-19',
        'user-123',
        'Digital Photography',
        'claude-4-sonnet',
        0,
        0,
        1749710060,
        1749710000,
        1749710060
    ),
    (
        'test-chat-20',
        'user-123',
        'Car Maintenance Basics',
        'qwen3',
        0,
        0,
        1749827819,
        1749760000,
        1749827819
    );

INSERT INTO
    messages (
        id,
        chat_id,
        role,
        model,
        content,
        reasoning,
        created_at,
        updated_at
    )
VALUES
    (
        'msg-11-1',
        'test-chat-11',
        'user',
        'claude-4-sonnet',
        'What are some simple ways to eat healthier?',
        '',
        1749309419,
        1749309419
    ),
    (
        'msg-11-2',
        'test-chat-11',
        'assistant',
        'claude-4-sonnet',
        'Start by incorporating more fruits and vegetables into your daily meals. Aim for a variety of colors to ensure a wide range of nutrients.',
        '',
        1749309439,
        1749309439
    ),
    (
        'msg-11-3',
        'test-chat-11',
        'user',
        'claude-4-sonnet',
        'Should I cut out carbs completely?',
        '',
        1749309459,
        1749309459
    ),
    (
        'msg-11-4',
        'test-chat-11',
        'assistant',
        'claude-4-sonnet',
        'Not necessarily! Focus on complex carbohydrates like whole grains, oats, and brown rice, which provide sustained energy and fiber, rather than eliminating them entirely.',
        '',
        1749309479,
        1749309479
    );

INSERT INTO
    messages (
        id,
        chat_id,
        role,
        model,
        content,
        reasoning,
        created_at,
        updated_at
    )
VALUES
    (
        'msg-12-1',
        'test-chat-12',
        'user',
        'qwen3',
        'What''s the best way to learn a new language quickly?',
        '',
        1749360000,
        1749360000
    ),
    (
        'msg-12-2',
        'test-chat-12',
        'assistant',
        'qwen3',
        'Immersion is key! Try to expose yourself to the language as much as possible through movies, music, and conversations with native speakers.',
        'The user is asking for the best way to learn a language quickly. I will focus on the principle of immersion, as it is widely considered the most effective for rapid acquisition. This includes various forms of exposure (media, native speakers) to maximize learning opportunities.',
        1749360020,
        1749360020
    ),
    (
        'msg-12-3',
        'test-chat-12',
        'user',
        'qwen3',
        'Are language learning apps effective?',
        '',
        1749360040,
        1749360040
    ),
    (
        'msg-12-4',
        'test-chat-12',
        'assistant',
        'qwen3',
        'Yes, apps like Duolingo or Babbel can be great for building vocabulary and basic grammar, but they are most effective when combined with other learning methods.',
        'The user is asking about the effectiveness of language learning apps. I will affirm their utility for specific foundational skills (vocabulary, grammar) while adding the nuance that they are best used as supplementary tools rather than standalone methods, suggesting a holistic approach.',
        1749360060,
        1749360060
    );

INSERT INTO
    messages (
        id,
        chat_id,
        role,
        model,
        content,
        reasoning,
        created_at,
        updated_at
    )
VALUES
    (
        'msg-13-1',
        'test-chat-13',
        'user',
        'claude-4-sonnet',
        'I want to redecorate my living room. Any budget-friendly ideas?',
        '',
        1749410000,
        1749410000
    ),
    (
        'msg-13-2',
        'test-chat-13',
        'assistant',
        'claude-4-sonnet',
        'Absolutely! Repainting walls, adding throw pillows and blankets, and rearranging existing furniture can make a big impact without spending much.',
        '',
        1749410020,
        1749410020
    ),
    (
        'msg-13-3',
        'test-chat-13',
        'user',
        'claude-4-sonnet',
        'What about lighting?',
        '',
        1749410040,
        1749410040
    ),
    (
        'msg-13-4',
        'test-chat-13',
        'assistant',
        'claude-4-sonnet',
        'Strategic lighting can transform a room. Consider adding floor lamps, table lamps, and even fairy lights to create different moods and highlight features.',
        '',
        1749410060,
        1749410060
    );

INSERT INTO
    messages (
        id,
        chat_id,
        role,
        model,
        content,
        reasoning,
        created_at,
        updated_at
    )
VALUES
    (
        'msg-14-1',
        'test-chat-14',
        'user',
        'qwen3',
        'I struggle with time management. How can I improve?',
        '',
        1749460000,
        1749460000
    ),
    (
        'msg-14-2',
        'test-chat-14',
        'assistant',
        'qwen3',
        'Try the Pomodoro Technique: work for 25 minutes, then take a 5-minute break. This helps maintain focus and prevents burnout.',
        'The user is asking for time management improvement. I will suggest a specific, well-known technique (Pomodoro) that directly addresses focus and burnout, as it provides a structured approach to work intervals and breaks.',
        1749460020,
        1749460020
    ),
    (
        'msg-14-3',
        'test-chat-14',
        'user',
        'qwen3',
        'How do I prioritize tasks?',
        '',
        1749460040,
        1749460040
    ),
    (
        'msg-14-4',
        'test-chat-14',
        'assistant',
        'qwen3',
        'Use the Eisenhower Matrix: categorize tasks as urgent/important, important/not urgent, urgent/not important, and neither. Focus on the urgent and important ones first.',
        'The user is asking for task prioritization methods. I will suggest the Eisenhower Matrix, a widely used and effective framework for categorizing and prioritizing tasks based on urgency and importance, allowing for clear decision-making on what to focus on first.',
        1749460060,
        1749460060
    );

INSERT INTO
    messages (
        id,
        chat_id,
        role,
        model,
        content,
        reasoning,
        created_at,
        updated_at
    )
VALUES
    (
        'msg-15-1',
        'test-chat-15',
        'user',
        'claude-4-sonnet',
        'My cat is scratching the furniture. Any advice?',
        '',
        1749510000,
        1749510000
    ),
    (
        'msg-15-2',
        'test-chat-15',
        'assistant',
        'claude-4-sonnet',
        'Provide plenty of scratching posts of various textures and heights. Encourage use with catnip or toys, and trim your cat''s nails regularly.',
        '',
        1749510020,
        1749510020
    ),
    (
        'msg-15-3',
        'test-chat-15',
        'user',
        'claude-4-sonnet',
        'How often should I feed my dog?',
        '',
        1749510040,
        1749510040
    ),
    (
        'msg-15-4',
        'test-chat-15',
        'assistant',
        'claude-4-sonnet',
        'Most adult dogs do well with two meals a day, morning and evening. Puppies may need 3-4 smaller meals daily.',
        '',
        1749510060,
        1749510060
    );

INSERT INTO
    messages (
        id,
        chat_id,
        role,
        model,
        content,
        reasoning,
        created_at,
        updated_at
    )
VALUES
    (
        'msg-16-1',
        'test-chat-16',
        'user',
        'qwen3',
        'What exactly is Artificial Intelligence?',
        '',
        1749560000,
        1749560000
    ),
    (
        'msg-16-2',
        'test-chat-16',
        'assistant',
        'qwen3',
        'Artificial Intelligence (AI) refers to the simulation of human intelligence in machines programmed to think like humans and mimic their actions, such as learning, problem-solving, and decision-making.',
        'The user is asking for a definition of Artificial Intelligence. I will provide a concise and clear explanation, focusing on the core concept of machines simulating human intelligence and key cognitive abilities like learning and problem-solving, which are central to AI.',
        1749560020,
        1749560020
    ),
    (
        'msg-16-3',
        'test-chat-16',
        'user',
        'qwen3',
        'What''s the difference between AI and Machine Learning?',
        '',
        1749560040,
        1749560040
    ),
    (
        'msg-16-4',
        'test-chat-16',
        'assistant',
        'qwen3',
        'Machine Learning is a subset of AI that enables systems to learn from data without explicit programming. AI is the broader concept of creating intelligent machines, while ML is a method to achieve that intelligence.',
        'The user is asking for the distinction between AI and Machine Learning. I will explain their relationship by defining Machine Learning as a subset of AI. This clarifies that ML is a specific approach or technique used within the broader field of AI to achieve intelligence through data-driven learning.',
        1749560060,
        1749560060
    );

INSERT INTO
    messages (
        id,
        chat_id,
        role,
        model,
        content,
        reasoning,
        created_at,
        updated_at
    )
VALUES
    (
        'msg-17-1',
        'test-chat-17',
        'user',
        'claude-4-sonnet',
        'Where can I find coding challenges for beginners?',
        '',
        1749610000,
        1749610000
    ),
    (
        'msg-17-2',
        'test-chat-17',
        'assistant',
        'claude-4-sonnet',
        'Websites like HackerRank, LeetCode (filter by easy), and freeCodeCamp offer a wide range of coding challenges perfect for beginners to practice their skills.',
        '',
        1749610020,
        1749610020
    ),
    (
        'msg-17-3',
        'test-chat-17',
        'user',
        'claude-4-sonnet',
        'What''s a good first project idea?',
        '',
        1749610040,
        1749610040
    ),
    (
        'msg-17-4',
        'test-chat-17',
        'assistant',
        'claude-4-sonnet',
        'A simple to-do list application, a basic calculator, or a rock-paper-scissors game are excellent first projects to solidify your foundational coding concepts.',
        '',
        1749610060,
        1749610060
    );

INSERT INTO
    messages (
        id,
        chat_id,
        role,
        model,
        content,
        reasoning,
        created_at,
        updated_at
    )
VALUES
    (
        'msg-18-1',
        'test-chat-18',
        'user',
        'qwen3',
        'How can I start practicing mindfulness?',
        '',
        1749660000,
        1749660000
    ),
    (
        'msg-18-2',
        'test-chat-18',
        'assistant',
        'qwen3',
        'Begin with simple breathing exercises. Find a quiet spot, sit comfortably, and focus on your breath for 5-10 minutes. Notice the sensation of air entering and leaving your body.',
        'The user is asking how to start mindfulness. I will suggest a foundational and accessible practice: focusing on breath. This approach is widely taught as a primary entry point into mindfulness because it requires no special equipment or prior experience, and can be done anywhere.',
        1749660020,
        1749660020
    ),
    (
        'msg-18-3',
        'test-chat-18',
        'user',
        'qwen3',
        'Is meditation difficult?',
        '',
        1749660040,
        1749660040
    ),
    (
        'msg-18-4',
        'test-chat-18',
        'assistant',
        'qwen3',
        'It can feel challenging at first, but consistency is more important than perfection. Don''t worry if your mind wanders; gently bring your focus back to your breath. Guided meditations can be very helpful.',
        'The user is asking about the difficulty of meditation. I will acknowledge initial challenges but emphasize consistency over perfection, reassuring them that mind-wandering is normal. I will also suggest guided meditations as a practical tool to assist beginners.',
        1749660060,
        1749660060
    );

INSERT INTO
    messages (
        id,
        chat_id,
        role,
        model,
        content,
        reasoning,
        created_at,
        updated_at
    )
VALUES
    (
        'msg-19-1',
        'test-chat-19',
        'user',
        'claude-4-sonnet',
        'What are the basics of digital photography?',
        '',
        1749710000,
        1749710000
    ),
    (
        'msg-19-2',
        'test-chat-19',
        'assistant',
        'claude-4-sonnet',
        'The "exposure triangle" is fundamental: aperture (controls depth of field), shutter speed (controls motion blur), and ISO (controls sensitivity to light). Understanding how they interact is key.',
        '',
        1749710020,
        1749710020
    ),
    (
        'msg-19-3',
        'test-chat-19',
        'user',
        'claude-4-sonnet',
        'How important is composition?',
        '',
        1749710040,
        1749710040
    ),
    (
        'msg-19-4',
        'test-chat-19',
        'assistant',
        'claude-4-sonnet',
        'Composition is extremely important! It''s how you arrange elements within your frame. Rules like the Rule of Thirds, leading lines, and framing can dramatically improve your photos.',
        '',
        1749710060,
        1749710060
    );

INSERT INTO
    messages (
        id,
        chat_id,
        role,
        model,
        content,
        reasoning,
        created_at,
        updated_at
    )
VALUES
    (
        'msg-20-1',
        'test-chat-20',
        'user',
        'qwen3',
        'What basic car maintenance should I know?',
        '',
        1749760000,
        1749760000
    ),
    (
        'msg-20-2',
        'test-chat-20',
        'assistant',
        'qwen3',
        'Regular oil changes, checking tire pressure, and ensuring your fluid levels (coolant, brake fluid, windshield washer fluid) are adequate are crucial for car longevity and safety.',
        'The user is asking for basic car maintenance. I will list essential, routine checks that directly impact the car''s longevity and safety. These include oil changes, tire pressure, and fluid levels, as they are fundamental to preventing common issues and ensuring safe operation.',
        1749760020,
        1749760020
    ),
    (
        'msg-20-3',
        'test-chat-20',
        'user',
        'qwen3',
        'How often should I check my oil?',
        '',
        1749760040,
        1749760040
    ),
    (
        'msg-20-4',
        'test-chat-20',
        'assistant',
        'qwen3',
        'It''s a good idea to check your engine oil level once a month, and definitely before a long road trip. This helps catch potential issues early.',
        'The user is asking for the frequency of oil checks. I will provide practical guidelines: monthly routine checks for general maintenance and an additional check before long trips, emphasizing that this proactive approach helps in early detection of potential issues.',
        1749827819,
        1749827819
    );

INSERT INTO
    attachments (id, message_id, name, type)
VALUES
    (
        'att-1-1-1',
        'msg-1-1',
        'morning_routine_inspiration.jpeg',
        'image/jpeg'
    ),
    (
        'att-1-3-1',
        'msg-1-3',
        'morning_routine_tracker.pdf',
        'application/pdf'
    ),
    (
        'att-3-1-1',
        'msg-3-1',
        'garden_layout_sketch.png',
        'image/png'
    ),
    (
        'att-3-3-1',
        'msg-3-3',
        'gardening_tool_list.docx',
        'application/vnd.openxmlformats-officedocument.wordprocessingml.document'
    ),
    (
        'att-5-1-1',
        'msg-5-1',
        'current_expenses_spreadsheet.xlsx',
        'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'
    ),
    (
        'att-5-3-1',
        'msg-5-3',
        'budget_template.pdf',
        'application/pdf'
    ),
    (
        'att-7-1-1',
        'msg-7-1',
        'leaky_faucet_photo.jpeg',
        'image/jpeg'
    ),
    (
        'att-7-3-1',
        'msg-7-3',
        'diy_repair_manual_cover.png',
        'image/png'
    ),
    (
        'att-9-1-1',
        'msg-9-1',
        'investment_portfolio_chart.svg',
        'image/svg+xml'
    ),
    (
        'att-9-3-1',
        'msg-9-3',
        'market_outlook_report.pdf',
        'application/pdf'
    ),
    (
        'att-11-1-1',
        'msg-11-1',
        'meal_plan_example.pdf',
        'application/pdf'
    ),
    (
        'att-11-3-1',
        'msg-11-3',
        'dietary_guidelines.pdf',
        'application/pdf'
    ),
    (
        'att-13-1-1',
        'msg-13-1',
        'living_room_before_photo.jpeg',
        'image/jpeg'
    ),
    (
        'att-13-3-1',
        'msg-13-3',
        'lighting_fixture_options.png',
        'image/png'
    ),
    (
        'att-15-1-1',
        'msg-15-1',
        'cat_scratch_post_ideas.jpeg',
        'image/jpeg'
    ),
    (
        'att-15-3-1',
        'msg-15-3',
        'dog_feeding_chart.pdf',
        'application/pdf'
    ),
    (
        'att-17-1-1',
        'msg-17-1',
        'coding_challenge_platform_screenshot.png',
        'image/png'
    ),
    (
        'att-17-3-1',
        'msg-17-3',
        'project_idea_flowchart.svg',
        'image/svg+xml'
    ),
    (
        'att-19-1-1',
        'msg-19-1',
        'camera_settings_diagram.jpeg',
        'image/jpeg'
    ),
    (
        'att-19-3-1',
        'msg-19-3',
        'composition_rule_of_thirds.png',
        'image/png'
    );