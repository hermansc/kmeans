package main

import (
    "math/rand"
    "math"
    "time"
    "github.com/ajstarks/svgo"
    "fmt"
    "bytes"
    "flag"
    "strconv"
    "net/http"
)

type Circle struct {
    x float64
    y float64
    cluster int
}

const (
    WIDTH = 70
    HEIGHT = 70
)

func CreateNormCircles(num int, devx, devy, meanx, meany float64) []*Circle {
    // Create a set of normally distributed points in the 2D space.
    // Parameters are given as flags or in form-data.
    circles := make([]*Circle, num)
    for i := 0; i < num; i++ {
        circles[i] = &Circle{
            x: (rand.NormFloat64() * devx + meanx),
            y: (rand.NormFloat64() * devy + meany),
        }
    }
    return circles
}

func CreateCircles(num int) []*Circle {
    // Create a set of evenly distributed circles in the 2D space.
    circles := make([]*Circle, num)
    for i := 0; i < num; i++ {
        circles[i] = &Circle{
            x: (rand.Float64() * WIDTH),
            y: (rand.Float64() * HEIGHT),
        }
    }
    return circles
}

func CalculateClusters(points []*Circle, centroids []Circle) {
    // For each point, we calculate distance to centroids.
    for i := 0; i < len(points); i++ {
        min := math.MaxFloat64
        for j := 0; j < len(centroids); j++ {
            p := points[i]
            c := centroids[j]

            // Find distance with euclidian distance. This can easily be
            // changed to another distance measure.
            dist := _euclidianDistance(p.x,c.x,p.y,c.y)

            if dist < min {
                // We have found a shorter distance to a centroid, and thus
                // change the cluster for this point.
                min = dist
                points[i].cluster = j
            }
        }
    }
}

func KMeans(points []*Circle, k, limit int) string {
    start := time.Now()

    if limit == 0 {
        // We probably don't more iterations than this anyway.
        limit = 10000
    }

    // Select starting points / centroids
    shuffeled := _shuffleElements(points)[:k]
    centroids := make([]Circle, len(shuffeled))
    for i := 0; i < len(shuffeled); i++ {
        centroids[i] = Circle{
            x: shuffeled[i].x,
            y: shuffeled[i].y,
            cluster: i,
        }
    }

    // For each centroid, calculate initial possition.
    CalculateClusters(points, centroids)

    change := true
    iterations := 0
    for (change && iterations < limit) {
        change = false // This will be changed until we converge

        for i := 0; i < len(centroids); i++ {
            // For each centroid, calculate new position based on average.
            sumX := 0.0
            sumY := 0.0
            n := 0
            for j := 0; j < len(points); j++ {
                if points[j].cluster == i {
                    sumX += points[j].x
                    sumY += points[j].y
                    n++
                }
            }

            // New X and Y coordinates are the average of all coordinates in
            // the cluster.
            newX := (sumX / float64(n))
            newY := (sumY / float64(n))

            // Check if these values are the same as before (in order to do
            // convergence checking
            if (newX != centroids[i].x) {
                centroids[i].x = newX
                change = true
            }
            if (newY != centroids[i].y) {
                centroids[i].y = newY
                change = true
            }
        }

        if (change) {
            // We've moved our centroids, so we recalculate cluster memberships
            CalculateClusters(points, centroids)
        }

        iterations++
    }

    // Calculate the milliseconds used
    used := float64(time.Now().Sub(start).Nanoseconds()) / 1000.0 / 1000.0

    // Print/return the solution
    s := fmt.Sprintf("<p>Found solution after <b>%v iterations (%.3f milliseconds)</b></p>", iterations, used)
    return s + OutputToSVG(points, centroids, k)
}

func OutputToSVG(points []*Circle, centroids []Circle, k int) string {
    // We output as a SVG, printing to a byte buffer
    var buffer bytes.Buffer
    canvas := svg.New(&buffer)

    // Stupid SVG API does not take floats, so we scale everything
    scale := 10.0
    canvas.Start(WIDTH * int(scale), HEIGHT * int(scale))

    // Get some (hopefully) nice colors for each different cluster
    colors := _randomColors(k)

    // Output the clusters and points
    for i := 0; i < len(points); i++ {
        color := colors[points[i].cluster]
        x := int(points[i].x * scale)
        y := int(points[i].y * scale)
        canvas.Circle(x, y, 3, "fill:" + color)
    }

    // Output the centroids in special colors (as a bit larger circles)
    for i := 0; i < len(centroids); i++ {
        color := colors[centroids[i].cluster]
        x := int(centroids[i].x * scale)
        y := int(centroids[i].y * scale)
        canvas.Circle(x, y, 6, "stroke:black;fill:" + color)
    }

    canvas.End()
    return buffer.String()
}

func kMeansHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "<h1>K-means demo</h1>")
    fmt.Fprintf(w,
    `<p>A tool for creating a 2D exhaustive partitional clustering using K-means and euclidian distance</p>
    <form action='.' method='get'>
        <b>K</b>: <input type='text' name='k'></input><br/>
        <b>Points</b>: <input type='text' name='points'></input></br>
        <b>Max iterations</b> (0 defaults to infinite): <input type='text' name='limit' value='0'></input></br>
        <b>Normal distribution?</b> Leave these blank if you want an even
        distribution.</br />
        <b>Deviation X</b>: <input type='text' name='devx'></input>
        <b>Deviation Y</b>: <input type='text' name='devy'></input></br>
        <b>Mean X</b>: <input type='text' name='meanx'></input>
        <b>Mean Y</b>: <input type='text' name='meany'></input></br>
        <button type='submit'>Visualize!</button>
    </form>`)

    k_s := r.FormValue("k")
    points_s := r.FormValue("points")
    limit_s := r.FormValue("limit")
    devx_s := r.FormValue("devx")
    meanx_s := r.FormValue("meanx")
    devy_s := r.FormValue("devy")
    meany_s := r.FormValue("meany")

    // Convert form-data to ints and floats.
    k, _ := strconv.Atoi(k_s)
    points, _ := strconv.Atoi(points_s)
    limit, _ := strconv.Atoi(limit_s)
    devx, _ := strconv.ParseFloat(devx_s, 64)
    meanx, _ := strconv.ParseFloat(meanx_s, 64)
    devy, _ := strconv.ParseFloat(devy_s, 64)
    meany, _ := strconv.ParseFloat(meany_s, 64)

    // First ensure that some needed varibles are set
    if (k == 0 || points == 0) {
        return
    }

    // Ensure k <= number of points
    if (k > points) {
        fmt.Fprintf(w, "<span style='color:red'><b>K needs to be <= number of points</b></span>")
        return
    }

    var circles []*Circle
    if (devx != 0.0 && meanx != 0.0 && devy != 0.0 && meany != 0.0) {
        // A bit ugly, but we check if all the deviation and mean variables are
        // set. Could in theory have 0 as mean, but this is ignored for now.
        circles = CreateNormCircles(points, devx, devy, meanx, meany)
    } else {
        // We create circles using an even distribution.
        circles = CreateCircles(points)
    }

    // Just ensure that k and points are set in order to do K-means.
    str := KMeans(circles, k, limit)
    fmt.Fprintf(w, str)
}

func main() {
    rand.Seed(time.Now().UTC().UnixNano())

    // Flags available to the binary
    var points = flag.Int("points", 100, "How many random points?")
    var k = flag.Int("k", 5, "Which value for K?")
    var lim = flag.Int("lim", 0, "Limit the number of iterations")
    var httpservice = flag.Bool("http", false, "Run as HTTP service?")

    var devx = flag.Float64("devx", 0.0, "X Deviation for calculating normal distribution")
    var devy = flag.Float64("devy", 0.0, "Y Deviation for calculating normal distribution")
    var meanx = flag.Float64("meanx", 0.0, "X Mean for calculating normal distribution")
    var meany = flag.Float64("meany", 0.0, "Y Deviation for calculating normal distribution")
    flag.Parse()

    if (*httpservice) {
        // If we run as web-service all parameters are handled in the form.
        http.HandleFunc("/", kMeansHandler)
        http.ListenAndServe(":8080", nil)
    } else  {
        if (*k == 0 || *points == 0) {
            // We can't calculate k-means when these values are 0.
            return
        }

        // Create the specified number of circles.
        var circles []*Circle
        if (*devy != 0 && *devx != 0 && *meanx != 0 && *meany != 0) {
            circles = CreateNormCircles(*points, *devx, *devy, *meanx, *meany)
        } else {
            circles = CreateCircles(*points)
        }

        // Run the algorithm on these circles, until we reach a convergence.
        str := KMeans(circles, *k, *lim)

        // Output the actual SVG / HTML data to stdout.
        fmt.Printf("%v\n", str)
    }
}

/*
    Helper functions
*/

func _randomColors(k int) []string {
    colors := []string{
        "#ef4444",
        "#faa31b",
        "#009f75",
        "#fff000",
        "#82c341",
        "#88c6ed",
        "#394ba0",
        "#d54799",
    }
    if k <= 8 {
        return colors[:k]
    }
    randomColors := make([]string, k)
    for i := 0; i < k; i++ {
        randomColors[i] = "#" + _randomHex()
    }
    return randomColors
}

func _randomHex() string {
    return fmt.Sprintf("%x", rand.Intn(16777215))
}

func _euclidianDistance(x1,x2,y1,y2 float64) float64 {
    // 2d euclidian distance, but could easily be extended for more dimensions.
    p1 := x1 - x2
    p2 := y1 - y2
    return math.Sqrt((p1*p1) + (p2*p2))
}

func _shuffleElements(array []*Circle) []*Circle {
    // Fisher-Yates shuffle.
    i := len(array) - 1
    for (i != 0) {
        j := rand.Intn(i + 1)
        tmp := array[i]
        array[i] = array[j]
        array[j] = tmp
        i--
    }
    return array
}
