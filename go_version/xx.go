package main 
import (
	"encoding/csv"
	"fmt"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"time"
)
func k_means_clust_new(data_list map[int][]float64, num_clust int, num_iter int, w int) ([][]float64, map[int][]int, map[int]float64) {

	// generate init centorids
	var keys []int
	for k, _ := range data_list {
		keys = append(keys, k)
	}
	rand_keys := (generateRandomNumber(0, len(keys), num_clust))
	var centroids [][]float64
	for _, v := range rand_keys {
		centroids = append(centroids, data_list[keys[v]])
	}
	//fmt.Println(rand_keys)
	//fmt.Println(data_list)
	// fmt.Println("xxxxx", rand_keys,centroids)

	counter := 0
	assignments := make(map[int][]int)
	sumdistance := make(map[int]float64)
	for i := 0; i < num_iter; i++ {
		counter += 1
		// fmt.Println("counter times", counter, "centroids", centroids)
		// init empty  assignment every iteration
		for k, _ := range centroids {
			assignments[k] = []int{}
			sumdistance[k] = 0.0
		}
		//增加每个质心的误差值保存
		for k, v := range data_list {
			// fmt.Println("pppppppppppppppp")
			min_dist := math.Inf(1)
			var closest_clust int
			// fmt.Println(closest_clust)
			for kk, vv := range centroids {
				// fmt.Println(LB_Keogh(v,vv,w))
				// fmt.Println(v, vv, w)
				if LB_Keogh(v, vv, w) < min_dist {
					cur_dist := DtwDistance(v, vv)
					// fmt.Println(cur_dist)
					if cur_dist < min_dist {
						min_dist = cur_dist
						closest_clust = kk
					}
				}
			}
			// fmt.Println(closest_clust)
			assignments[closest_clust] = append(assignments[closest_clust], k)
			sumdistance[closest_clust] += min_dist
			// fmt.Println("iter times",i,"assignment",assignments)
		}

		for k, v := range assignments {
			sumdistance[k] = sumdistance[k] / float64(len(v))
			var clust_sum []float64
			for _, vv := range v {
				for kkk, vvv := range data_list[vv] {
					if len(clust_sum) < kkk+1 {
						clust_sum = append(clust_sum, 0)
					}
					clust_sum[kkk] += vvv
				}
			}
			for kk, vv := range clust_sum {
				centroids[k][kk] = vv / float64(len(v))
			}

		}
	}
	return centroids, assignments, sumdistance
}

func bisecting_k_means_clust(data_list map[int][]float64, num_clust int, num_iter int, w int) ([][]float64, [][]int) {
	//var clust_sum []float64
	// generate init centorids
	var allcentroids [][]float64
	var allassign [][]int
	var alldistance []float64
	for {
		//fmt.Println("xxxxxxxxxxxxx")
		//fmt.Println(len(allcentroids))
		var lastd = math.Inf(1)
		lastassign := make(map[int][]int)
		lastdistance := make(map[int]float64)
		var lastcentroids [][]float64
		//求最小值
		//fmt.Println("yyyyyyyyyyyyyyyy")
		for i := 0; i < 20; i++ {
		// for i := 0; i < 100; i++ {

			//fmt.Println("zzz")
			centroids, assignments, sumdistance := k_means_clust_new(data_list, 2, num_iter, w)
			//fmt.Println("iii")
			nowdistance := sum([]float64{sumdistance[0], sumdistance[1]})
			if nowdistance < lastd {
				lastd = nowdistance
				lastcentroids = centroids
				lastassign = assignments
				lastdistance = sumdistance
			}
		}
		// 汇总
		//fmt.Println(lastcentroids)
		//fmt.Println(lastassign)
		//fmt.Println(lastdistance)
		for k, v := range lastcentroids {
			allcentroids = append(allcentroids, v)
			allassign = append(allassign, lastassign[k])
			alldistance = append(alldistance, lastdistance[k])
		}
		//fmt.Println(allcentroids)
		//fmt.Println(allassign)
		//fmt.Println(alldistance)
		if len(allcentroids) == num_clust {
			break
		}
		// 求最大值
		var maxdistance = math.Inf(-1)
		var max_index = 0

		for k, v := range alldistance {
			if v > maxdistance {
				maxdistance = v
				max_index = k
			}
		}
		fmt.Println(max_index)
		max_index_list := allassign[max_index]
		fmt.Println(max_index_list)
		tmp := make(map[int][]float64)
		for _, v := range max_index_list {
			tmp[v] = data_list[v]
		}
		//fmt.Println(tmp)
		data_list = tmp
		//fmt.Println(data_list)
		allcentroids = DeleteSlice(allcentroids, max_index)
		//fmt.Println(allcentroids)
		allassign = DeleteSlice3(allassign, max_index)
		//fmt.Println(allassign)
		alldistance = DeleteSlice2(alldistance, max_index)
		//fmt.Println(alldistance)
		//fmt.Println(DeleteSlice3(allcentroids, max_index), e)

		//DeleteSlice(allassign,max_index)
		//DeleteSlice(,max_index)
	}

	fmt.Println(allassign)
	return allcentroids, allassign
}

//删除切片
func DeleteSlice(sss [][]float64, index int) ([][]float64) {
	//sliceValue := reflect.ValueOf(slice)
	length := len(sss)
	if sss == nil || length == 0 || (length-1) < index {
		return nil
	}
	if length-1 == index {
		return sss[0:index]
	} else if (length - 1) >= index {
		//return reflect.AppendSlice(sliceValue.Slice(0, index), sliceValue.Slice(index+1, length)).Interface()
		tmp := sss[0: index]
		for _, v := range sss[index+1:length] {
			tmp = append(tmp, v)
		}
		return tmp
	}
	return nil
}
func DeleteSlice2(sss []float64, index int) ([]float64) {
	//sliceValue := reflect.ValueOf(slice)
	length := len(sss)
	if sss == nil || length == 0 || (length-1) < index {
		return nil
	}
	if length-1 == index {
		return sss[0:index]
	} else if (length - 1) >= index {
		//return reflect.AppendSlice(sliceValue.Slice(0, index), sliceValue.Slice(index+1, length)).Interface()
		tmp := sss[0: index]
		for _, v := range sss[index+1:length] {
			tmp = append(tmp, v)
		}
		return tmp
	}
	return nil
}

func DeleteSlice3(sss [][]int, index int) ([][]int) {
	//sliceValue := reflect.ValueOf(slice)
	length := len(sss)
	if sss == nil || length == 0 || (length-1) < index {
		return nil
	}
	if length-1 == index {
		return sss[0:index]
	} else if (length - 1) >= index {
		//return reflect.AppendSlice(sliceValue.Slice(0, index), sliceValue.Slice(index+1, length)).Interface()
		tmp := sss[0: index]
		for _, v := range sss[index+1:length] {
			tmp = append(tmp, v)
		}
		return tmp
	}
	return nil
}

func ShortData(data map[string][]float64,start int,end int)(map[string][]float64){
	shortData := make(map[string][]float64)
	for i := range data{
		// fmt.Println(data[i][start:end])
		shortData[i] = data[i][start:end]
	}
	return shortData
}

func readcsv(path string) map[string][]float64{
	// read csv part
	file, err := os.Open(path)
	if err != nil {
		// err is printable
		// elements passed are separated by space automatically
		fmt.Println("error:", err)
	}
	// automatically call Close() at the end of current method
	defer file.Close()
	reader := csv.NewReader(file)
	reader.Comma = ','
	rawCSVdata, err := reader.ReadAll()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// rawCSVdata = rawCSVdata[:20]

	// sanity check, display to standard output
	// for _, each := range rawCSVdata {
	// fmt.Printf("email : %s and timestamp : %s\n", each[0], each[1])
	// }
	newdata := dataclean(rawCSVdata)
	return newdata
}

func dataclean(raw [][]string) map[string][]float64 {
	csv_data := make(map[string][]float64)
	for _, line := range raw {
		for _, nums := range line {
			if n, err := strconv.ParseFloat(nums, 64); err == nil {
				csv_data[line[0]] = append([]float64{n}, csv_data[line[0]]...)
			}
		}
	}
	return csv_data
}

func get_centroid(datas map[string][]float64,n int) ([][]float64, map[int][]int,[]string,map[string][]float64) {
	var keys = sorted_keys(datas)
	// fmt.Println(keys)
	var data_list [][]float64
	// var data_map map[string][]float64
	data_map := make(map[string][]float64)
	for _, v := range keys {
		a := to_zero(datas[v])
		data_list = append(data_list, a)
		data_map[v] = a
	}
	// fmt.Println(data_list)
	centroids, assignments := k_means_clust(data_list, n, 20, 3)
	// fmt.Println("okkkkkkkkkkk")
	// fmt.Println(centroids)
	// fmt.Println(assignments)
	return centroids, assignments,keys,data_map
}
func get_centroid_new(datas map[string][]float64,n int) ([][]float64, map[int][]int,[]string,map[string][]float64) {
	var keys = sorted_keys(datas)
	// fmt.Println(keys)
	// var data_list map[int][]float64
	data_list := make(map[int][]float64)
	// var data_map map[string][]float64
	data_map := make(map[string][]float64)
	for k, v := range keys {
		a := to_zero(datas[v])
		data_list[k] = a
		data_map[v] = a
	}
	// fmt.Println(data_list)
	centroids, assignments := bisecting_k_means_clust(data_list, n, 20, 3)
	// fmt.Println("okkkkkkkkkkk")
	// fmt.Println(centroids)
	// fmt.Println(assignments)
	newassignments := make(map[int][]int)
	for k,v := range assignments{
		newassignments[k] = v
	}
	fmt.Println(newassignments,keys,data_map)
	return centroids, newassignments,keys,data_map
}


func to_zero(arr []float64) []float64 {
	var tmp []float64
	for _, v := range arr {
		tmp = append(tmp, Round(v-arr[0], 2))
	}
	// fmt.Println(tmp)
	return tmp
}

func to_rate(arr []float64) []float64 {
	var tmp []float64
	for k, v := range arr {
		// tmp = append(tmp, Round(v-arr[0], 2))
		// tmp = append(tmp, Round(v/arr[0], 2))
		if k!=0{
		tmp = append(tmp,Round((v-arr[k-1])/arr[k-1],2))
		}
	}
	// fmt.Println(tmp)
	return tmp
}

func sorted_keys(m map[string][]float64) []string {
	sorted_keys := make([]string, 0)
	for k, _ := range m {
		sorted_keys = append(sorted_keys, k)
	}
	sort.Strings(sorted_keys)
	return sorted_keys
}

func Round(val float64, places int) float64 {
	var t float64
	f := math.Pow10(places)
	x := val * f
	if math.IsInf(x, 0) || math.IsNaN(x) {
		return val
	}
	if x >= 0.0 {
		t = math.Ceil(x)
		if (t - x) > 0.50000000001 {
			t -= 1.0
		}
	} else {
		t = math.Ceil(-x)
		if (t + x) > 0.50000000001 {
			t -= 1.0
		}
		t = -t
	}
	x = t / f

	if !math.IsInf(x, 0) {
		return x
	}

	return t
}

func generateRandomNumber(start int, end int, count int) []int {
	//范围检查
	if end < start || (end-start) < count {
		return nil
	}
	//存放结果的slice
	nums := make([]int, 0)
	//随机数生成器，加入时间戳保证每次生成的随机数不一样
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(nums) < count {
		//生成随机数
		num := r.Intn((end - start)) + start
		//查重
		exist := false
		for _, v := range nums {
			if v == num {
				exist = true
				break
			}
		}
		if !exist {
			nums = append(nums, num)
		}
	}
	return nums
}

func max(data []float64) float64 {
	maximum := math.Inf(-1)
	for i := 0; i < len(data); i++ {
		maximum = math.Max(data[i], maximum)
	}
	return maximum
}

func min(data []float64) float64 {
	minimum := math.Inf(1)
	for i := 0; i < len(data); i++ {
		minimum = math.Min(data[i], minimum)
	}
	return minimum
}

func sum(data []float64) float64 {
	all := 0.0
	for i := 0; i < len(data); i++ {
		all += data[i]
	}
	return all
}

func DtwDistance(s1 []float64, s2 []float64) float64 {
	DTW := make(map[[2]int]float64)
	for i := -1; i < len(s1); i++ {
		for j := -1; j < len(s2); j++ {
			keyarr := [2]int{i, j}
			DTW[keyarr] = math.Inf(1)
		}
	}
	DTW[[2]int{-1, -1}] = 0

	w := math.Max(5, math.Abs(float64(len(s1)-len(s2))))
	for i := 0; i < len(s1); i++ {
		lower := int(math.Max(float64(0), float64(i)-w))
		upper := int(math.Min(float64(len(s2)), float64(i)+w))
		for j := lower; j < upper; j++ {
			dist := math.Pow(s1[i]-s2[j], 2)
			// fmt.Println(i,j,s1[i]+s2[j],dist)
			values := []float64{DTW[[2]int{i, j - 1}], DTW[[2]int{i - 1, j}], DTW[[2]int{i - 1, j - 1}]}
			// DTW[[2]int{i + 1, j + 1}] = dist + min(values)
			DTW[[2]int{i, j}] = dist + min(values)
			// fmt.Println(i,j,dist,values,min(values))
			// fmt.Println(DTW[[2]int{i, j}])
		}
	}
	return math.Sqrt(DTW[[2]int{len(s1) - 1, len(s2) - 1}])
}

func LB_Keogh(s1 []float64, s2 []float64, r int) float64 {
	// fmt.Println(len(s1),len(s2),r)
	LB_sum := 0.0
	for ind, i := range s1 {
		start := 0
		if ind-r >= 0 {
			start = ind - r
		}
		// fmt.Println(start,ind+r)
		end := ind + r
		if end >= len(s1) {
			end = len(s1)
		}
		lower_bound := (min(s2[start:end]))
		upper_bound := (max(s2[start:end]))

		if i > upper_bound {
			LB_sum = LB_sum + math.Pow((i-upper_bound), 2)
		} else if i < lower_bound {
			LB_sum = LB_sum + math.Pow((i-lower_bound), 2)
		}
	}
	return LB_sum
}
func k_means_clust(data_list [][]float64, num_clust int, num_iter int, w int) ([][]float64, map[int][]int) {

	// generate init centorids
	rand_keys := (generateRandomNumber(0, len(data_list), num_clust))
	var centroids [][]float64
	for _, v := range rand_keys {
		centroids = append(centroids, data_list[v])
	}
	// fmt.Println("xxxxx", rand_keys,centroids)

	counter := 0
	assignments := make(map[int][]int)
	for i := 0; i < num_iter; i++ {
		counter += 1
		// fmt.Println("counter times", counter, "centroids", centroids)
		// init empty  assignment every iteration
		for k, _ := range centroids {
			assignments[k] = []int{}
		}
		for k, v := range data_list {
			// fmt.Println("pppppppppppppppp")
			min_dist := math.Inf(1)
			var closest_clust int
			// fmt.Println(closest_clust)
			for kk, vv := range centroids {
				// fmt.Println(LB_Keogh(v,vv,w))
				if LB_Keogh(v, vv, w) < min_dist {
					cur_dist := DtwDistance(v, vv)
					// fmt.Println(cur_dist)
					if cur_dist < min_dist {
						min_dist = cur_dist
						closest_clust = kk
					}
				}
			}
			// fmt.Println(closest_clust)
			assignments[closest_clust] = append(assignments[closest_clust], k)
			// fmt.Println("iter times",i,"assignment",assignments)
		}

		for k, v := range assignments {
			var clust_sum []float64
			for _, vv := range v {
				for kkk, vvv := range data_list[vv] {
					if len(clust_sum) < kkk+1 {
						clust_sum = append(clust_sum, 0)
					}
					clust_sum[kkk] += vvv
				}
			}
			for kk, vv := range clust_sum {
				centroids[k][kk] = vv / float64(len(v))
			}

		}
	}
	return centroids, assignments
}

func get_stock_map(stocklist []int)(map[int]int){
	stockmap := make(map[int]int)
	for k,v := range stocklist{
		stockmap[v] = k
	}
	return stockmap
}

var stocklist = []float64{
	20110104,
	20110105,
	20110106,
	20110107,
	20110110,
	20110111,
	20110112,
	20110113,
	20110114,
	20110117,
	20110118,
	20110119,
	20110120,
	20110121,
	20110124,
	20110125,
	20110126,
	20110127,
	20110128,
	20110131,
	20110201,
	20110209,
	20110210,
	20110211,
	20110214,
	20110215,
	20110216,
	20110217,
	20110218,
	20110221,
	20110222,
	20110223,
	20110224,
	20110225,
	20110228,
	20110301,
	20110302,
	20110303,
	20110304,
	20110307,
	20110308,
	20110309,
	20110310,
	20110311,
	20110314,
	20110315,
	20110316,
	20110317,
	20110318,
	20110321,
	20110322,
	20110323,
	20110324,
	20110325,
	20110328,
	20110329,
	20110330,
	20110331,
	20110401,
	20110406,
	20110407,
	20110408,
	20110411,
	20110412,
	20110413,
	20110414,
	20110415,
	20110418,
	20110419,
	20110420,
	20110421,
	20110422,
	20110425,
	20110426,
	20110427,
	20110428,
	20110429,
	20110503,
	20110504,
	20110505,
	20110506,
	20110509,
	20110510,
	20110511,
	20110512,
	20110513,
	20110516,
	20110517,
	20110518,
	20110519,
	20110520,
	20110523,
	20110524,
	20110525,
	20110526,
	20110527,
	20110530,
	20110531,
	20110601,
	20110602,
	20110603,
	20110607,
	20110608,
	20110609,
	20110610,
	20110613,
	20110614,
	20110615,
	20110616,
	20110617,
	20110620,
	20110621,
	20110622,
	20110623,
	20110624,
	20110627,
	20110628,
	20110629,
	20110630,
	20110701,
	20110704,
	20110705,
	20110706,
	20110707,
	20110708,
	20110711,
	20110712,
	20110713,
	20110714,
	20110715,
	20110718,
	20110719,
	20110720,
	20110721,
	20110722,
	20110725,
	20110726,
	20110727,
	20110728,
	20110729,
	20110801,
	20110802,
	20110803,
	20110804,
	20110805,
	20110808,
	20110809,
	20110810,
	20110811,
	20110812,
	20110815,
	20110816,
	20110817,
	20110818,
	20110819,
	20110822,
	20110823,
	20110824,
	20110825,
	20110826,
	20110829,
	20110830,
	20110831,
	20110901,
	20110902,
	20110905,
	20110906,
	20110907,
	20110908,
	20110909,
	20110913,
	20110914,
	20110915,
	20110916,
	20110919,
	20110920,
	20110921,
	20110922,
	20110923,
	20110926,
	20110927,
	20110928,
	20110929,
	20110930,
	20111010,
	20111011,
	20111012,
	20111013,
	20111014,
	20111017,
	20111018,
	20111019,
	20111020,
	20111021,
	20111024,
	20111025,
	20111026,
	20111027,
	20111028,
	20111031,
	20111101,
	20111102,
	20111103,
	20111104,
	20111107,
	20111108,
	20111109,
	20111110,
	20111111,
	20111114,
	20111115,
	20111116,
	20111117,
	20111118,
	20111121,
	20111122,
	20111123,
	20111124,
	20111125,
	20111128,
	20111129,
	20111130,
	20111201,
	20111202,
	20111205,
	20111206,
	20111207,
	20111208,
	20111209,
	20111212,
	20111213,
	20111214,
	20111215,
	20111216,
	20111219,
	20111220,
	20111221,
	20111222,
	20111223,
	20111226,
	20111227,
	20111228,
	20111229,
	20111230,
	20120104,
	20120105,
	20120106,
	20120109,
	20120110,
	20120111,
	20120112,
	20120113,
	20120116,
	20120117,
	20120118,
	20120119,
	20120120,
	20120130,
	20120131,
	20120201,
	20120202,
	20120203,
	20120206,
	20120207,
	20120208,
	20120209,
	20120210,
	20120213,
	20120214,
	20120215,
	20120216,
	20120217,
	20120220,
	20120221,
	20120222,
	20120223,
	20120224,
	20120227,
	20120228,
	20120229,
	20120301,
	20120302,
	20120305,
	20120306,
	20120307,
	20120308,
	20120309,
	20120312,
	20120313,
	20120314,
	20120315,
	20120316,
	20120319,
	20120320,
	20120321,
	20120322,
	20120323,
	20120326,
	20120327,
	20120328,
	20120329,
	20120330,
	20120405,
	20120406,
	20120409,
	20120410,
	20120411,
	20120412,
	20120413,
	20120416,
	20120417,
	20120418,
	20120419,
	20120420,
	20120423,
	20120424,
	20120425,
	20120426,
	20120427,
	20120502,
	20120503,
	20120504,
	20120507,
	20120508,
	20120509,
	20120510,
	20120511,
	20120514,
	20120515,
	20120516,
	20120517,
	20120518,
	20120521,
	20120522,
	20120523,
	20120524,
	20120525,
	20120528,
	20120529,
	20120530,
	20120531,
	20120601,
	20120604,
	20120605,
	20120606,
	20120607,
	20120608,
	20120611,
	20120612,
	20120613,
	20120614,
	20120615,
	20120618,
	20120619,
	20120620,
	20120621,
	20120625,
	20120626,
	20120627,
	20120628,
	20120629,
	20120702,
	20120703,
	20120704,
	20120705,
	20120706,
	20120709,
	20120710,
	20120711,
	20120712,
	20120713,
	20120716,
	20120717,
	20120718,
	20120719,
	20120720,
	20120723,
	20120724,
	20120725,
	20120726,
	20120727,
	20120730,
	20120731,
	20120801,
	20120802,
	20120803,
	20120806,
	20120807,
	20120808,
	20120809,
	20120810,
	20120813,
	20120814,
	20120815,
	20120816,
	20120817,
	20120820,
	20120821,
	20120822,
	20120823,
	20120824,
	20120827,
	20120828,
	20120829,
	20120830,
	20120831,
	20120903,
	20120904,
	20120905,
	20120906,
	20120907,
	20120910,
	20120911,
	20120912,
	20120913,
	20120914,
	20120917,
	20120918,
	20120919,
	20120920,
	20120921,
	20120924,
	20120925,
	20120926,
	20120927,
	20120928,
	20121008,
	20121009,
	20121010,
	20121011,
	20121012,
	20121015,
	20121016,
	20121017,
	20121018,
	20121019,
	20121022,
	20121023,
	20121024,
	20121025,
	20121026,
	20121029,
	20121030,
	20121031,
	20121101,
	20121102,
	20121105,
	20121106,
	20121107,
	20121108,
	20121109,
	20121112,
	20121113,
	20121114,
	20121115,
	20121116,
	20121119,
	20121120,
	20121121,
	20121122,
	20121123,
	20121126,
	20121127,
	20121128,
	20121129,
	20121130,
	20121203,
	20121204,
	20121205,
	20121206,
	20121207,
	20121210,
	20121211,
	20121212,
	20121213,
	20121214,
	20121217,
	20121218,
	20121219,
	20121220,
	20121221,
	20121224,
	20121225,
	20121226,
	20121227,
	20121228,
	20121231,
	20130104,
	20130107,
	20130108,
	20130109,
	20130110,
	20130111,
	20130114,
	20130115,
	20130116,
	20130117,
	20130118,
	20130121,
	20130122,
	20130123,
	20130124,
	20130125,
	20130128,
	20130129,
	20130130,
	20130131,
	20130201,
	20130204,
	20130205,
	20130206,
	20130207,
	20130208,
	20130218,
	20130219,
	20130220,
	20130221,
	20130222,
	20130225,
	20130226,
	20130227,
	20130228,
	20130301,
	20130304,
	20130305,
	20130306,
	20130307,
	20130308,
	20130311,
	20130312,
	20130313,
	20130314,
	20130315,
	20130318,
	20130319,
	20130320,
	20130321,
	20130322,
	20130325,
	20130326,
	20130327,
	20130328,
	20130329,
	20130401,
	20130402,
	20130403,
	20130408,
	20130409,
	20130410,
	20130411,
	20130412,
	20130415,
	20130416,
	20130417,
	20130418,
	20130419,
	20130422,
	20130423,
	20130424,
	20130425,
	20130426,
	20130502,
	20130503,
	20130506,
	20130507,
	20130508,
	20130509,
	20130510,
	20130513,
	20130514,
	20130515,
	20130516,
	20130517,
	20130520,
	20130521,
	20130522,
	20130523,
	20130524,
	20130527,
	20130528,
	20130529,
	20130530,
	20130531,
	20130603,
	20130604,
	20130605,
	20130606,
	20130607,
	20130613,
	20130614,
	20130617,
	20130618,
	20130619,
	20130620,
	20130621,
	20130624,
	20130625,
	20130626,
	20130627,
	20130628,
	20130701,
	20130702,
	20130703,
	20130704,
	20130705,
	20130708,
	20130709,
	20130710,
	20130711,
	20130712,
	20130715,
	20130716,
	20130717,
	20130718,
	20130719,
	20130722,
	20130723,
	20130724,
	20130725,
	20130726,
	20130729,
	20130730,
	20130731,
	20130801,
	20130802,
	20130805,
	20130806,
	20130807,
	20130808,
	20130809,
	20130812,
	20130813,
	20130814,
	20130815,
	20130816,
	20130819,
	20130820,
	20130821,
	20130822,
	20130823,
	20130826,
	20130827,
	20130828,
	20130829,
	20130830,
	20130902,
	20130903,
	20130904,
	20130905,
	20130906,
	20130909,
	20130910,
	20130911,
	20130912,
	20130913,
	20130916,
	20130917,
	20130918,
	20130923,
	20130924,
	20130925,
	20130926,
	20130927,
	20130930,
	20131008,
	20131009,
	20131010,
	20131011,
	20131014,
	20131015,
	20131016,
	20131017,
	20131018,
	20131021,
	20131022,
	20131023,
	20131024,
	20131025,
	20131028,
	20131029,
	20131030,
	20131031,
	20131101,
	20131104,
	20131105,
	20131106,
	20131107,
	20131108,
	20131111,
	20131112,
	20131113,
	20131114,
	20131115,
	20131118,
	20131119,
	20131120,
	20131121,
	20131122,
	20131125,
	20131126,
	20131127,
	20131128,
	20131129,
	20131202,
	20131203,
	20131204,
	20131205,
	20131206,
	20131209,
	20131210,
	20131211,
	20131212,
	20131213,
	20131216,
	20131217,
	20131218,
	20131219,
	20131220,
	20131223,
	20131224,
	20131225,
	20131226,
	20131227,
	20131230,
	20131231,
	20140102,
	20140103,
	20140106,
	20140107,
	20140108,
	20140109,
	20140110,
	20140113,
	20140114,
	20140115,
	20140116,
	20140117,
	20140120,
	20140121,
	20140122,
	20140123,
	20140124,
	20140127,
	20140128,
	20140129,
	20140130,
	20140207,
	20140210,
	20140211,
	20140212,
	20140213,
	20140214,
	20140217,
	20140218,
	20140219,
	20140220,
	20140221,
	20140224,
	20140225,
	20140226,
	20140227,
	20140228,
	20140303,
	20140304,
	20140305,
	20140306,
	20140307,
	20140310,
	20140311,
	20140312,
	20140313,
	20140314,
	20140317,
	20140318,
	20140319,
	20140320,
	20140321,
	20140324,
	20140325,
	20140326,
	20140327,
	20140328,
	20140331,
	20140401,
	20140402,
	20140403,
	20140404,
	20140408,
	20140409,
	20140410,
	20140411,
	20140414,
	20140415,
	20140416,
	20140417,
	20140418,
	20140421,
	20140422,
	20140423,
	20140424,
	20140425,
	20140428,
	20140429,
	20140430,
	20140505,
	20140506,
	20140507,
	20140508,
	20140509,
	20140512,
	20140513,
	20140514,
	20140515,
	20140516,
	20140519,
	20140520,
	20140521,
	20140522,
	20140523,
	20140526,
	20140527,
	20140528,
	20140529,
	20140530,
	20140603,
	20140604,
	20140605,
	20140606,
	20140609,
	20140610,
	20140611,
	20140612,
	20140613,
	20140616,
	20140617,
	20140618,
	20140619,
	20140620,
	20140623,
	20140624,
	20140625,
	20140626,
	20140627,
	20140630,
	20140701,
	20140702,
	20140703,
	20140704,
	20140707,
	20140708,
	20140709,
	20140710,
	20140711,
	20140714,
	20140715,
	20140716,
	20140717,
	20140718,
	20140721,
	20140722,
	20140723,
	20140724,
	20140725,
	20140728,
	20140729,
	20140730,
	20140731,
	20140801,
	20140804,
	20140805,
	20140806,
	20140807,
	20140808,
	20140811,
	20140812,
	20140813,
	20140814,
	20140815,
	20140818,
	20140819,
	20140820,
	20140821,
	20140822,
	20140825,
	20140826,
	20140827,
	20140828,
	20140829,
	20140901,
	20140902,
	20140903,
	20140904,
	20140905,
	20140909,
	20140910,
	20140911,
	20140912,
	20140915,
	20140916,
	20140917,
	20140918,
	20140919,
	20140922,
	20140923,
	20140924,
	20140925,
	20140926,
	20140929,
	20140930,
	20141008,
	20141009,
	20141010,
	20141013,
	20141014,
	20141015,
	20141016,
	20141017,
	20141020,
	20141021,
	20141022,
	20141023,
	20141024,
	20141027,
	20141028,
	20141029,
	20141030,
	20141031,
	20141103,
	20141104,
	20141105,
	20141106,
	20141107,
	20141110,
	20141111,
	20141112,
	20141113,
	20141114,
	20141117,
	20141118,
	20141119,
	20141120,
	20141121,
	20141124,
	20141125,
	20141126,
	20141127,
	20141128,
	20141201,
	20141202,
	20141203,
	20141204,
	20141205,
	20141208,
	20141209,
	20141210,
	20141211,
	20141212,
	20141215,
	20141216,
	20141217,
	20141218,
	20141219,
	20141222,
	20141223,
	20141224,
	20141225,
	20141226,
	20141229,
	20141230,
	20141231,
	20150105,
	20150106,
	20150107,
	20150108,
	20150109,
	20150112,
	20150113,
	20150114,
	20150115,
	20150116,
	20150119,
	20150120,
	20150121,
	20150122,
	20150123,
	20150126,
	20150127,
	20150128,
	20150129,
	20150130,
	20150202,
	20150203,
	20150204,
	20150205,
	20150206,
	20150209,
	20150210,
	20150211,
	20150212,
	20150213,
	20150216,
	20150217,
	20150225,
	20150226,
	20150227,
	20150302,
	20150303,
	20150304,
	20150305,
	20150306,
	20150309,
	20150310,
	20150311,
	20150312,
	20150313,
	20150316,
	20150317,
	20150318,
	20150319,
	20150320,
	20150323,
	20150324,
	20150325,
	20150326,
	20150327,
	20150330,
	20150331,
	20150401,
	20150402,
	20150403,
	20150407,
	20150408,
	20150409,
	20150410,
	20150413,
	20150414,
	20150415,
	20150416,
	20150417,
	20150420,
	20150421,
	20150422,
	20150423,
	20150424,
	20150427,
	20150428,
	20150429,
	20150430,
	20150504,
	20150505,
	20150506,
	20150507,
	20150508,
	20150511,
	20150512,
	20150513,
	20150514,
	20150515,
	20150518,
	20150519,
	20150520,
	20150521,
	20150522,
	20150525,
	20150526,
	20150527,
	20150528,
	20150529,
	20150601,
	20150602,
	20150603,
	20150604,
	20150605,
	20150608,
	20150609,
	20150610,
	20150611,
	20150612,
	20150615,
	20150616,
	20150617,
	20150618,
	20150619,
	20150623,
	20150624,
	20150625,
	20150626,
	20150629,
	20150630,
	20150701,
	20150702,
	20150703,
	20150706,
	20150707,
	20150708,
	20150709,
	20150710,
	20150713,
	20150714,
	20150715,
	20150716,
	20150717,
	20150720,
	20150721,
	20150722,
	20150723,
	20150724,
	20150727,
	20150728,
	20150729,
	20150730,
	20150731,
	20150803,
	20150804,
	20150805,
	20150806,
	20150807,
	20150810,
	20150811,
	20150812,
	20150813,
	20150814,
	20150817,
	20150818,
	20150819,
	20150820,
	20150821,
	20150824,
	20150825,
	20150826,
	20150827,
	20150828,
	20150831,
	20150901,
	20150902,
	20150907,
	20150908,
	20150909,
	20150910,
	20150911,
	20150914,
	20150915,
	20150916,
	20150917,
	20150918,
	20150921,
	20150922,
	20150923,
	20150924,
	20150925,
	20150928,
	20150929,
	20150930,
	20151008,
	20151009,
	20151012,
	20151013,
	20151014,
	20151015,
	20151016,
	20151019,
	20151020,
	20151021,
	20151022,
	20151023,
	20151026,
	20151027,
	20151028,
	20151029,
	20151030,
	20151102,
	20151103,
	20151104,
	20151105,
	20151106,
	20151109,
	20151110,
	20151111,
	20151112,
	20151113,
	20151116,
	20151117,
	20151118,
	20151119,
	20151120,
	20151123,
	20151124,
	20151125,
	20151126,
	20151127,
	20151130,
	20151201,
	20151202,
	20151203,
	20151204,
	20151207,
	20151208,
	20151209,
	20151210,
	20151211,
	20151214,
	20151215,
	20151216,
	20151217,
	20151218,
	20151221,
	20151222,
	20151223,
	20151224,
	20151225,
	20151228,
	20151229,
	20151230,
	20151231,
	20160104,
	20160105,
	20160106,
	20160107,
	20160108,
	20160111,
	20160112,
	20160113,
	20160114,
	20160115,
	20160118,
	20160119,
	20160120,
	20160121,
	20160122,
	20160125,
	20160126,
	20160127,
	20160128,
	20160129,
	20160201,
	20160202,
	20160203,
	20160204,
	20160205,
	20160215,
	20160216,
	20160217,
	20160218,
	20160219,
	20160222,
	20160223,
	20160224,
	20160225,
	20160226,
	20160229,
	20160301,
	20160302,
	20160303,
	20160304,
	20160307,
	20160308,
	20160309,
	20160310,
	20160311,
	20160314,
	20160315,
	20160316,
	20160317,
	20160318,
	20160321,
	20160322,
	20160323,
	20160324,
	20160325,
	20160328,
	20160329,
	20160330,
	20160331,
	20160401,
	20160405,
	20160406,
	20160407,
	20160408,
	20160411,
	20160412,
	20160413,
	20160414,
	20160415,
	20160418,
	20160419,
	20160420,
	20160421,
	20160422,
	20160425,
	20160426,
	20160427,
	20160428,
	20160429,
	20160503,
	20160504,
	20160505,
	20160506,
	20160509,
	20160510,
	20160511,
	20160512,
	20160513,
	20160516,
	20160517,
	20160518,
	20160519,
	20160520,
	20160523,
	20160524,
	20160525,
	20160526,
	20160527,
	20160530,
	20160531,
	20160601,
	20160602,
	20160603,
	20160606,
	20160607,
	20160608,
	20160613,
	20160614,
	20160615,
	20160616,
	20160617,
	20160620,
	20160621,
	20160622,
	20160623,
	20160624,
	20160627,
	20160628,
	20160629,
	20160630,
	20160701,
	20160704,
	20160705,
	20160706,
	20160707,
	20160708,
	20160711,
	20160712,
	20160713,
	20160714,
	20160715,
	20160718,
	20160719,
	20160720,
	20160721,
	20160722,
	20160725,
	20160726,
	20160727,
	20160728,
	20160729,
	20160801,
	20160802,
	20160803,
	20160804,
	20160805,
	20160808,
	20160809,
	20160810,
	20160811,
	20160812,
	20160815,
	20160816,
	20160817,
	20160818,
	20160819,
	20160822,
	20160823,
	20160824,
	20160825,
	20160826,
	20160829,
	20160830,
	20160831,
	20160901,
	20160902,
	20160905,
	20160906,
	20160907,
	20160908,
	20160909,
	20160912,
	20160913,
	20160914,
	20160919,
	20160920,
	20160921,
	20160922,
	20160923,
	20160926,
	20160927,
	20160928,
	20160929,
	20160930,
	20161010,
	20161011,
	20161012,
	20161013,
	20161014,
	20161017,
	20161018,
	20161019,
	20161020,
	20161021,
	20161024,
	20161025,
	20161026,
	20161027,
	20161028,
	20161031,
	20161101,
	20161102,
	20161103,
	20161104,
	20161107,
	20161108,
	20161109,
	20161110,
	20161111,
	20161114,
	20161115,
	20161116,
	20161117,
	20161118,
	20161121,
	20161122,
	20161123,
	20161124,
	20161125,
	20161128,
	20161129,
	20161130,
	20161201,
	20161202,
	20161205,
	20161206,
	20161207,
	20161208,
	20161209,
	20161212,
	20161213,
	20161214,
	20161215,
	20161216,
	20161219,
	20161220,
	20161221,
	20161222,
	20161223,
	20161226,
	20161227,
	20161228,
	20161229,
	20161230,
	20170103,
	20170104,
	20170105,
	20170106,
	20170109,
	20170110,
	20170111,
	20170112,
	20170113,
	20170116,
	20170117,
	20170118,
	20170119,
	20170120,
	20170123,
	20170124,
	20170125,
	20170126,
	20170203,
	20170206,
	20170207,
	20170208,
	20170209,
	20170210,
	20170213,
	20170214,
	20170215,
	20170216,
	20170217,
	20170220,
	20170221,
	20170222,
	20170223,
	20170224,
	20170227,
	20170228,
	20170301,
	20170302,
	20170303,
	20170306,
	20170307,
	20170308,
	20170309,
	20170310,
	20170313,
	20170314,
	20170315,
	20170316,
	20170317,
	20170320,
	20170321,
	20170322,
	20170323,
	20170324,
	20170327,
	20170328,
	20170329,
	20170330,
	20170331,
	20170405,
	20170406,
	20170407,
	20170410,
	20170411,
	20170412,
	20170413,
	20170414,
	20170417,
	20170418,
	20170419,
	20170420,
	20170421,
	20170424,
	20170425,
	20170426,
	20170427,
	20170428,
	20170502,
	20170503,
	20170504,
	20170505,
	20170508,
	20170509,
	20170510,
	20170511,
	20170512,
	20170515,
	20170516,
	20170517,
	20170518,
	20170519,
	20170522,
	20170523,
	20170524,
	20170525,
	20170526,
	20170531,
	20170601,
	20170602,
	20170605,
	20170606,
	20170607,
	20170608,
	20170609,
	20170612,
	20170613,
	20170614,
	20170615,
	20170616,
	20170619,
	20170620,
	20170621,
	20170622,
	20170623,
	20170626,
	20170627,
	20170628,
	20170629,
	20170630,
	20170703,
	20170704,
	20170705,
	20170706,
	20170707,
	20170710,
	20170711,
	20170712,
	20170713,
	20170714,
	20170717,
	20170718,
	20170719,
	20170720,
	20170721,
	20170724,
	20170725,
	20170726,
	20170727,
	20170728,
	20170731,
	20170801,
}