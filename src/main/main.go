package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"net/http"
)

//Структура для входных параметров
//Имена переменныех должны начинаться с большой буквы
type Data struct {
	X float64 `json: "x"`
	Y float64 `json: "y"`
}

//Структура данных результата
type RData struct {
	Phii float64 `json: "r"`
}

//////// Начало функции вычисления /////////////////////////////////////////////////////////

// Функция вычисления Тау

func tauCount(minusX []float64, minusY []float64, plusX []float64, plusY []float64, elAmount int, k, dphi float64, r float64) float64 {

	var minusSum, plusSum, sums float64 = 0, 0, 0
	var tau float64

	for i := 0; i < elAmount/2; i++ {
		minusSum += (math.Log((math.Sqrt(math.Pow(plusX[elAmount/4]-r-minusX[i], 2) + math.Pow(plusY[elAmount/4]-minusY[i], 2))) / (math.Sqrt(math.Pow(minusX[i], 2) + math.Pow(minusY[i], 2)))))
		plusSum += (math.Log((math.Sqrt(math.Pow(plusX[elAmount/4]-r-plusX[i], 2) + math.Pow(plusY[elAmount/4]-plusY[i], 2))) / (math.Sqrt(math.Pow(plusX[i], 2) + math.Pow(plusY[i], 2)))))
	}
	sums = sums + minusSum - plusSum
	minusSum = 0
	plusSum = 0

	for i := 0; i < elAmount/2; i++ {
		minusSum = minusSum + (math.Log((math.Sqrt(math.Pow(plusX[elAmount/4]+r-minusX[i], 2) + math.Pow(plusY[elAmount/4]-minusY[i], 2))) / (math.Sqrt(math.Pow(minusX[i], 2) + math.Pow(minusY[i], 2)))))
		plusSum = plusSum + (math.Log((math.Sqrt(math.Pow(plusX[elAmount/4]+r-plusX[i], 2) + math.Pow(plusY[elAmount/4]-plusY[i], 2))) / (math.Sqrt(math.Pow(plusX[i], 2) + math.Pow(plusY[i], 2)))))
	}
	sums = sums + minusSum - plusSum
	minusSum = 0
	plusSum = 0

	for i := 0; i < elAmount/2; i++ {
		minusSum += math.Log((math.Sqrt(math.Pow(minusX[elAmount/4]-r-minusX[i], 2) + math.Pow(minusY[elAmount/4]-minusY[i], 2))) / (math.Sqrt(math.Pow(minusX[i], 2) + math.Pow(minusY[i], 2))))
		plusSum += math.Log((math.Sqrt(math.Pow(minusX[elAmount/4]-r-plusX[i], 2) + math.Pow(minusY[elAmount/4]-plusY[i], 2))) / (math.Sqrt(math.Pow(plusX[i], 2) + math.Pow(plusY[i], 2))))
	}
	sums = sums - minusSum + plusSum
	minusSum = 0
	plusSum = 0

	for i := 0; i < elAmount/2; i++ {
		minusSum += math.Log((math.Sqrt(math.Pow(minusX[elAmount/4]+r-minusX[i], 2) + math.Pow(minusY[elAmount/4]-minusY[i], 2))) / (math.Sqrt(math.Pow(minusX[i], 2) + math.Pow(minusY[i], 2))))
		plusSum += math.Log((math.Sqrt(math.Pow(minusX[elAmount/4]+r-plusX[i], 2) + math.Pow(minusY[elAmount/4]-plusY[i], 2))) / (math.Sqrt(math.Pow(plusX[i], 2) + math.Pow(plusY[i], 2))))
	}
	sums = sums - minusSum + plusSum
	sums = sums * k

	tau = dphi / sums

	return tau
}

// Функция вычисления ПФИ
func phiCount(x, y float64, minusX []float64, minusY []float64, plusX []float64, plusY []float64, elAmount int, k, tau float64) float64 {
	var minusSum, plusSum, phi float64 = 0, 0, 0

	for i := 0; i < elAmount/2; i++ {
		minusSum = minusSum + (math.Log((math.Sqrt(math.Pow(x-minusX[i], 2) + math.Pow(y-minusY[i], 2))) / (math.Sqrt(math.Pow(minusX[i], 2) + math.Pow(minusY[i], 2)))))
		plusSum = plusSum + (math.Log((math.Sqrt(math.Pow(x-plusX[i], 2) + math.Pow(y-plusY[i], 2))) / (math.Sqrt(math.Pow(plusX[i], 2) + math.Pow(plusY[i], 2)))))
	}
	phi = (2 * k * tau * minusSum) - (2 * k * tau * plusSum)
	return phi
}

// Функция округления
func Round(x, unit float64) float64 {
	return math.Round(x*100) / 100
}

////////////////////////////////////////////////////////////

//обработчик для функции Ajax запроса расчет потенциала поля по координатам
func myfuncAjax(w http.ResponseWriter, r *http.Request) {
	var d Data
	var rp RData
	//fmt.Println("%v",r.Body)
	//Читаем тело запроса r.Body и преобразуем данные в структуру Data
	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	///////////////////////////////////////////////////////////////

	var otrAmount float64 = 4 // Кол-во отражений
	var elAmount int          // Кол- элементов
	var stolb, strok int      // Кол-во строк столбцов
	var s, c, b, rr float64 = 11, 14.4, 9.9, 1.5
	//var length, height float64 = 41.8, 28.8
	var x, y float64
	var tau float64
	var phi, dphi float64 = 0, 14.5 // dhi- Максимальное пфи
	var k float64
	k = 1.0 / (4.0 * 3.14 * 8.85 * math.Pow(10, -12))
	var phiShift float64
	//default_random_engine generator;
	var otrAmountint int = 4
	for i := 1; i <= otrAmountint; i++ { // Вычисление кол-во элементов
		elAmount = elAmount + i

	}
	elAmount = elAmount * 16
	elAmount = elAmount + 2
	strok = (2 * otrAmountint) + 1 // ВЫчисление строк

	stolb = elAmount / strok // вычисление столбцов

	minusX := make([]float64, 81) // создание массивов координат
	minusY := make([]float64, 81)
	plusX := make([]float64, 81)
	plusY := make([]float64, 81)

	var jFloat float64 // Счетчики
	var j1Float float64

	var floati float64 = 1 // для 148, 149 строчки

	if otrAmountint%2 == 0 { // Если количество отражений кратно 2, подсчёт координат первой строки
		minusX[0] = -s - 2*b*otrAmount - 2*s*(otrAmount-1)
		plusX[0] = -s - 2*b*otrAmount - 2*s*otrAmount
		minusY[0] = 2 * c * otrAmount
		plusY[0] = 2 * c * otrAmount

		for j := 1; j < stolb/2; j++ {
			if j%2 == 0 {
				jFloat = 0
			} else {
				jFloat = 1
			}
			if (j+1)%2 == 0 {
				j1Float = 0
			} else {
				j1Float = 1
			}
			minusX[j] = minusX[j-1] + 2*b + s*4*j1Float //(j%2)
			plusX[j] = plusX[j-1] + 2*b + s*4*jFloat    //((j+1)%2)
			minusY[j] = minusY[0]
			plusY[j] = plusY[0]

		}

	}

	if otrAmountint%2 == 1 {
		minusX[0] = -s - 2*b*otrAmount - 2*s*otrAmount
		plusX[0] = -s - 2*b*otrAmount - 2*s*(otrAmount-1)
		minusY[0] = 2 * c * otrAmount
		plusY[0] = 2 * c * otrAmount
		for j := 1; j < stolb/2; j++ {
			if j%2 == 0 {
				jFloat = 0
			} else {
				jFloat = 1
			}
			if (j+1)%2 == 0 {
				j1Float = 0
			} else {
				j1Float = 1
			}
			minusX[j] = minusX[j-1] + 2*b + s*4*jFloat //(j%2)
			plusX[j] = plusX[j-1] + 2*b + s*4*j1Float  //((j+1)%2)
			minusY[j] = minusY[0]
			plusY[j] = plusY[0]
		}
	}

	floati = 0

	for i := 1; i < strok; i++ {
		floati = floati + 1
		for j := 0; j < stolb/2; j++ {
			minusX[j+i*stolb/2] = minusX[j]
			plusX[j+i*stolb/2] = plusX[j]
			minusY[j+i*stolb/2] = minusY[0] - (2 * c * floati)
			plusY[j+i*stolb/2] = plusY[0] - (2 * c * floati)
		}
	}

	tau = tauCount(minusX, minusY, plusX, plusY, elAmount, k, dphi, rr)

	phiShift = (phiCount(minusX[elAmount/4]-rr, minusY[elAmount/4], minusX, minusY, plusX, plusY, elAmount, k, tau) + phiCount(minusX[elAmount/4]+rr, minusY[elAmount/4], minusX, minusY, plusX, plusY, elAmount, k, tau)) / 2 // подсчёт сдвига потенциала
	for i := 1; i < 2; i++ {
		phi = 0

		x = Round(d.X*0.02645833333333-600.9448818898*0.02645833333333, 0.05)
		y = Round(d.Y*0.02645833333333-336.3779527559*0.02645833333333, 0.05)

		if ((math.Pow(x-plusX[elAmount/4], 2) + math.Pow(y-plusY[elAmount/4], 2)) <= math.Pow(rr, 2)) || ((math.Pow(x-minusX[elAmount/4], 2) + math.Pow(y-minusY[elAmount/4], 2)) <= math.Pow(rr, 2)) { // если точка в электроде

			if (math.Pow(x-plusX[elAmount/4], 2) + math.Pow(y-plusY[elAmount/4], 2)) <= math.Pow(rr, 2) { // если в плюсе

				phi = (phiCount(plusX[elAmount/4]-rr, plusY[elAmount/4], minusX, minusY, plusX, plusY, elAmount, k, tau) + phiCount(plusX[elAmount/4]+rr, plusY[elAmount/4], minusX, minusY, plusX, plusY, elAmount, k, tau)) / 2
			}
			if (math.Pow(x-minusX[elAmount/4], 2) + math.Pow(y-minusY[elAmount/4], 2)) <= math.Pow(rr, 2) { // если в минусе

				phi = (phiCount(minusX[elAmount/4]-rr, minusY[elAmount/4], minusX, minusY, plusX, plusY, elAmount, k, tau) + phiCount(minusX[elAmount/4]+rr, minusY[elAmount/4], minusX, minusY, plusX, plusY, elAmount, k, tau)) / 2
			}
			phi = phi - phiShift

		} else {
			phi = phiCount(x, y, minusX, minusY, plusX, plusY, elAmount, k, tau)
			phi = phi - phiShift
			//normal_distribution<double> distribution(phi, sigma); //
			//phi = distribution(generator);
			//cout << "phi = " << phi << endl;
			//fmt.Println("2 phi=")
			//fmt.Printf("%-10.3f", phi)
		}
	}

	////////Конец функции вычисления потенциала//////////////////////////////////////////////
	phi = Round(phi, 0.01)
	rp.Phii = phi
	fmt.Print(phi)

	// Сознаем json response from struct RData
	a, err := json.Marshal(rp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// Задаем заголовок данным для возврата ответа
	w.Header().Set("Content-Type", "application/json")
	//Возвращаем результат в виде json структуры
	w.Write(a)
}

func main() {
	// Heroku прокидывает порт для приложения в переменную окружения PORT
	port := os.Getenv("PORT")

	http.Handle("/", http.FileServer(http.Dir("static")))
	// Определяем маршрут и обработчик для функции Ajax запроса
	http.HandleFunc("/myfunc", myfuncAjax)
	go http.ListenAndServe(":"+port, nil)
}

//log.Fatal()  // Спросить
