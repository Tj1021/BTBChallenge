# BTBChallenge
The .go file will create two JSON files.
1.) entries.JSON will be the logs of login attempts with the normalized format. 
2.) numLogs.JSON is used to keep track of the last known log, so that the no log is pulled twice.

the visualization data can be reclaculated by setting the variable getVisualizationData to true.
the visualization data is stored in visualizationData.txt

the HTML file (visualization.html) uses css and javascript to create a bar chart (with the help of CanvasJS).
The bar chart tallies the total number of login attempts on each day of the week.
