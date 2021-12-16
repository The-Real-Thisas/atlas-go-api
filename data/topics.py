import os
import json
from rich import print

with open("data.json") as f:
    data = json.load(f)

main = {}

# for subject in data.keys():
#     for topic in data[subject].keys():
#         topicsFound = {}
#         for question in data[subject][topic].keys():
#             if data[subject][topic][question]["topics"] in topicsFound:
#                 topicsFound[data[subject][topic][question]["topics"]] += 1
#             else:
#                 topicsFound[data[subject][topic][question]["topics"]] = 1
#         try:
#             highestTopicCount = max(list(topicsFound.values()))
#         except ValueError:
#             print(f"{subject} {topic} has no questions")
#             pass
#         topicFound = ""
#         for topic in topicsFound.keys():
#             if topicsFound[topic] == highestTopicCount:
#                 topicFound = topic
#         print(topicFound)

"""
For each subject, for each topic, for each question, find the topic with the highest count.

Final output:
{
    "subject": [{"topic": "topic"}]
}
Example:
{
    "physics": [{"physics-1": "physics-1"}, {"physics-2": "physics-2"}]
}
"""
for subject in data.keys():
    subjectTopics = {}
    for topic in data[subject].keys():
        topicsFound = {}
        for question in data[subject][topic].keys():
            if data[subject][topic][question]["topics"] in topicsFound:
                topicsFound[data[subject][topic][question]["topics"]] += 1
            else:
                topicsFound[data[subject][topic][question]["topics"]] = 1
        try:
            highestTopicCount = max(list(topicsFound.values()))
        except ValueError:
            print(f"{subject} {topic} has no questions")
            highestTopicCount = 0
        topicFound = ""
        if not highestTopicCount == 0:
            for topicF in topicsFound.keys():
                if topicsFound[topicF] == highestTopicCount:
                    topicFound = topicF
        else:
            topicFound = topic
        subjectTopics[topicFound] = topic
    main[subject] = subjectTopics

print(main)

data = {}

"""
{
    "Subject": {
        "topics": []
        "topicNames": []
    }
}
"""

for subject in main.keys():
    topics = []
    topicNames = []
    for topic in main[subject].keys():
        topicNames.append(topic)
        topics.append(main[subject][topic])
    data[subject] = {"topics": topics, "topicNames": topicNames}

with open("topics.json", "w") as f:
    json.dump(data, f)
