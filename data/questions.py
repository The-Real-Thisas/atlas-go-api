import json

with open('./data.json') as json_file:
    data = json.load(json_file)
    print(data)

questions = {"questions": []}

for subject in data.keys():
    for topic in data[subject].keys():
        for question in data[subject][topic].keys():
            questionData = {
                "subject": subject,
                "topic": topic,
                "id": data[subject][topic][question]["id"],
                "paper_number": data[subject][topic][question]["paper_number"],
                "topics": data[subject][topic][question]["topics"],
                "question_urls": data[subject][topic][question]["question_urls"],
                "answer_urls": data[subject][topic][question]["answer_urls"]
            }
            questions["questions"].append(questionData)

print(questions)

with open('./questions.json', 'w') as outfile:
    json.dump(questions, outfile)