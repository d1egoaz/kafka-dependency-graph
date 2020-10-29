run:
	go run main.go "/Users/diegoalvarez/Downloads/kafka_message_lost_report-relation_between_topics_and_apps-06ecdc5d22e9-2020-10-28-15-10-26.csv" > mygraph.dot
	dot -Tsvg mygraph.dot -o mygraph.svg

clean:
	rm mygraph.*
