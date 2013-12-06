# K-means implementation in Go.

Provided is my source code for a K-means algorithm demonstration, implemented and written in Go (the new concurrent language from Google).
The demonstration was part of a university course in information retrieval and shows the goal of K-means: create clusters of data where the average distance between the centroids and particles are minimal.

## K-means clustering: a short introduction

We want to partition N objects into K clusters. We need to define the K ourselves. The need to define such a K is actually a disadvantage using K-means and if one want to automate it one can cross-validate with a range of K values and find K-values that creates a "knee" in the residual sum of squares (more on this can be found elsewhere).

K-means is a flat-clustering method and thus has the disadvantages of being:

- Non-deterministic
- Having no relations between clusters
- Not (necessarily) finding the optimal solution, but a "good enough (a common theme in NP-hard problems and heuristic algorithms)
- Easy to obtain singleton clusters (that is outliers), if particles are very spread.


The remedy for many of these disadvantages are using something like an heuristic clustering algorithm. However this comes at the cost of a larger running time.

## K-means: how does it work?

We select K random points in our 2D-coordinate space and coin these K-points centroids. Then we calculate the distance between all particles/points and the centroids and assign each point to its nearest centroid. Particles sharing the same assignment of centroids are in the same cluster.

When all particles are assigned we calculate the average X- and Y-values in *each cluster* and give the resulting coordinates from this calculation as new position to the centroid.

The final step is then to again, re-calculate the distances between all particles and the centroids and assign each point anew. These steps (assignment and updating centroids) are repeated until we hit a convergence.

## A note on the code

It can be run in two ways:

    $ ./kmeans --http

or

    $ ./kmeans --points=5000 --k=5

In the first we start a web-service using the CGI-library offered by Go. The system should be available at http://localhost:8080/

In the second way we will output the HTML to the stdout. I recommend redirecting this to a file:

    $ ./kmeans --points=5000 --k=5 > mykmeans.html

## Example

Currently the web-service is available at the following URL (likely to change): http://cassarossa.samfundet.no:8080/
