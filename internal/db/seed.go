package db

import (
	"context"
	"fmt"
	"log"
	"math/rand"

	"github.com/yunsuk-jeung/social/internal/store"
)

func Seed(store store.Storage) {
	ctx := context.Background()

	users := generateUsers(100)
	for _, user := range users {
		if err := store.Users.Create(ctx, user); err != nil {
			log.Println("Error creating user: ", err)
			return
		}
	}
	posts := generatePosts(200, users)
	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Println("Error creating post: ", err)
			return
		}
	}

	comments := generateComments(500, users, posts)
	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			log.Println("Error creating comment: ", err)
			return
		}
	}
	log.Println("Seeding complete")
}

var usernames = []string{
	"alice", "bob", "charlie", "david", "eve", "frank", "grace", "heidi", "ivan", "judy",
	"karl", "louis", "mallory", "nancy", "oscar", "peggy", "quentin", "ruth", "sam", "trent",
	"ursula", "victor", "wendy", "xander", "yvonne", "zack", "amber", "brian", "cindy", "doug",
	"ella", "felix", "gina", "hank", "irene", "jack", "karen", "leo", "mona", "nate",
	"olga", "paul", "quincy", "rose", "steve", "tina", "uma", "vince", "wanda", "xenia",
}

var blogTitles = []string{
	"How to Start Fresh",
	"5 Tips for Better Focus",
	"Lessons from Failure",
	"Why Simplicity Wins",
	"Building Good Habits",
	"Finding Your Passion",
	"The Power of Routine",
	"Small Changes, Big Results",
	"Staying Motivated Daily",
	"Mindset Shifts That Help",
	"What I Learned This Year",
	"The Joy of Minimalism",
	"Overcoming Fear Step-by-Step",
	"Balancing Work and Life",
	"Creativity on Busy Days",
	"Making Time for Growth",
	"Secrets of Productive People",
	"Tiny Wins Every Day",
	"Letting Go of Perfection",
	"Trusting the Process",
}

var blogContents = []string{
	"Starting fresh means letting go of past mistakes and focusing on today. Small steps every day can lead to big changes over time.",
	"Eliminate distractions, set clear goals, take regular breaks, prioritize your tasks, and practice mindfulness to sharpen your focus.",
	"Failure isn't the end—it's feedback. Every setback teaches resilience, adaptability, and the value of persistence.",
	"Simple systems are easier to maintain, easier to scale, and easier to enjoy. Focus on what matters most and cut the rest.",
	"Good habits are built through consistency, not motivation. Start small, track progress, and reward yourself for sticking with it.",
	"Passion often follows effort. Explore new things, reflect on what energizes you, and stay curious until you find your spark.",
	"A strong routine frees up mental energy, builds momentum, and keeps you moving even when motivation dips.",
	"Tiny improvements, repeated daily, can transform your life. Focus on getting 1% better every day.",
	"Motivation fluctuates. Build systems that keep you going, even on days you don't feel inspired.",
	"Switch from fixed thinking to growth thinking. See challenges as opportunities to learn, not threats to your identity.",
	"Reflection brings clarity. Every year teaches lessons in patience, resilience, and the importance of gratitude.",
	"Owning less means stressing less. Minimalism isn't about deprivation; it's about making space for what truly matters.",
	"Fear shrinks when you face it. Break challenges into small steps and celebrate each tiny victory.",
	"Work-life balance is dynamic, not static. Set boundaries, prioritize downtime, and communicate clearly with others.",
	"You don't need hours to be creative. Capture small ideas, write quick notes, and trust that creativity builds over time.",
	"If you wait for 'free time', it won't come. Schedule growth activities like you schedule meetings.",
	"Productive people focus deeply, say no often, automate the boring, and rest deliberately to recharge.",
	"Small achievements compound. Celebrate tiny wins—they build momentum and reinforce positive habits.",
	"Perfectionism paralyzes progress. Done is better than perfect—ship it, learn, and improve.",
	"Progress isn't always visible day-to-day. Trust the habits you've built and keep moving forward, even when it's hard to see results.",
}

var blogTags = [][]string{
	{"restart", "motivation", "new goals"},
	{"deep work", "focus hacks", "attention"},
	{"failure lessons", "bounce back", "resilience"},
	{"minimal life", "simple living", "essentials"},
	{"daily habits", "habit building", "consistency"},
	{"finding purpose", "career paths", "self discovery"},
	{"morning routine", "daily rhythm", "rituals"},
	{"small steps", "big dreams", "daily action"},
	{"self drive", "inner strength", "stay motivated"},
	{"growth mindset", "mental shift", "overcoming limits"},
	{"year in review", "reflection", "personal stories"},
	{"declutter", "minimalism journey", "space and peace"},
	{"face your fears", "confidence", "mental toughness"},
	{"balance", "work life harmony", "self care"},
	{"creative bursts", "idea generation", "art life"},
	{"self improvement", "time investment", "life design"},
	{"efficiency", "work smart", "high output"},
	{"daily wins", "habit tracker", "small victories"},
	{"ship fast", "get things done", "progress over perfection"},
	{"long game", "patience", "trust yourself"},
}

var blogComments = []string{
	"This post really spoke to me!", "I love the simplicity of your approach!", "Can’t wait to implement these steps!",
	"Such a helpful post! I struggle with focus, but I’ll try these tips.", "These productivity tips will be my new mantra.", "Love the reminder to take breaks and focus.",
	"This was so empowering. I’ve failed before, but now I see it differently.", "I used to fear failure, but now I embrace it.", "This article gives me so much hope!",
	"I’ve been wanting to simplify my life! Thanks for the great advice.", "Minimalism has helped me feel so much more grounded.", "I’m starting a decluttering challenge this weekend!",
	"I’m definitely going to try tracking my habits. Thanks for the motivation!", "Consistency is key! I love this post.", "It’s good to know that small wins really do add up.",
	"I feel so lost when it comes to finding my passion, but this is reassuring.", "This post made me realize that passion comes with time.", "Great read! I’ll explore more and stay patient.",
	"Routine is everything! I feel so much more productive when I stick to one.", "I love your take on routines. They help me get through the day.", "I’ll be starting a morning routine tomorrow!",
	"I totally agree! Small steps lead to big change. Love this perspective.", "Tiny wins really are a game changer.", "I’m going to track my small victories from now on.",
	"Motivation can be tough, but I’m going to create a system to keep myself going.", "This post has made me rethink how I stay motivated.", "Love the practical advice on staying motivated daily.",
	"I’ve been working on a growth mindset for a while, and this was so encouraging.", "This is the reminder I needed — challenges help us grow!", "Growth mindset for the win! This really resonated.",
	"Looking back at the year, I’ve learned so much. This was a great reflection.", "I love your perspective on reflecting on failures and successes.", "I’m going to take more time to reflect after every year.",
	"Minimalism is truly life-changing. I’ve been on this journey for a few months.", "Love this perspective on living with less. It really frees up space in your life!", "Decluttering has helped me mentally and physically.",
	"Facing fears is tough, but it’s the only way to grow. Great article!", "I feel more confident after reading this. I’m going to start facing my fears head-on.", "This post was a real confidence booster.",
	"I really need to work on my work-life balance. Thanks for the tips!", "This post has some great advice for juggling everything.", "Boundaries are so important for mental health. I’ll start practicing that.",
	"Creativity isn’t easy, but I’ll try to embrace it more every day.", "I love your tips on nurturing creativity in busy times.", "This is a great reminder to make space for creative moments.",
	"I need to be more intentional about scheduling time for growth.", "Love this advice on making growth a priority. I’m going to plan accordingly!", "This made me realize how important it is to prioritize personal growth.",
	"Efficiency is everything! I’m going to start working smarter, not harder.", "Such great advice on how to stay productive without burning out.", "I’m going to try using some of your strategies to improve my efficiency.",
	"This is why small wins matter so much. Thanks for sharing this!", "I’m starting to track my little wins now. They feel so rewarding.", "Love how you emphasize progress, not perfection!",
	"Perfectionism has always held me back, but this is freeing. Thanks for sharing.", "Done is better than perfect. I need to remember that.", "I’m embracing imperfection now — I feel lighter already.",
	"This was a great reminder to trust the process. Things take time!", "I’m going to work on being more patient with myself.", "I needed this post today! Trusting the journey is everything.",
}

func generateUsers(num int) []*store.User {
	users := make([]*store.User, num)
	for i := 0; i < num; i++ {
		users[i] = &store.User{
			Username: usernames[i%len(usernames)] + fmt.Sprintf("%d", i),
			Email:    usernames[i%len(usernames)] + fmt.Sprintf("%d", i) + "@example.com",
			Password: "123456",
		}

	}
	return users
}

func generatePosts(num int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, num)
	for i := 0; i < num; i++ {
		user := users[rand.Intn(len(users))]
		posts[i] = &store.Post{
			UserID:  user.ID,
			Title:   blogTitles[i%len(blogTitles)] + fmt.Sprintf("%d", i),
			Content: blogContents[i%len(blogContents)] + fmt.Sprintf("%d", i),
			Tags:    blogTags[i%len(blogTags)],
		}

	}
	return posts
}

func generateComments(num int, users []*store.User, posts []*store.Post) []*store.Comment {
	cms := make([]*store.Comment, num)
	for i := 0; i < num; i++ {
		cms[i] = &store.Comment{
			PostID:  posts[rand.Intn(len(posts))].ID,
			UserID:  users[rand.Intn(len(users))].ID,
			Content: blogComments[rand.Intn(len(blogComments))],
		}
	}
	return cms
}
