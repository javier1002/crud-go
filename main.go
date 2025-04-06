// compilamos aqui el codigo
package main

import (
	"fmt"         
	"database/sql"
	"log"           //terminar
	"net/http"      //web
	"text/template" //plantilla
	_ "github.com/go-sql-driver/mysql" //driver de mysql,_carga la un driver
)

func conexionBD() (conexion *sql.DB) { //funcion para conectar a la base de datos
	Driver := "mysql"
	usuario := "root"
	password := ""
	nombreBD := "sistema"
	conexion, err := sql.Open(Driver, usuario+":"+password+"@tcp(127.0.0.1)/"+nombreBD)

 //conexion a la base de datos
	if err != nil {         //si hay error
		panic(err.Error())  //muestra el error
	}
	return conexion 		//retorna la conexion

}
//funcion para iniciar el servidor	

var plantillas = template.Must(template.ParseGlob("plantillas/*")) //plantilla El archivo archivos.tmpl es un archivo de plantilla HTML que se usa en Go para generar páginas web dinámicas. La extensión
func main() {
	http.HandleFunc("/", Inicio)                 //llama la funcion inicio
	http.HandleFunc("/crear", Crear)             //recibe a la funcion crear de inicio y la envia a el metodo crear
	http.HandleFunc("/insertar", Insertar)       //recibe a la funcion crear de inicio y la envia a el metodo crear
	http.HandleFunc("/borrar", Borrar)           //recibe a la funcion crear de inicio y la envia a el metodo crear
	http.HandleFunc("/editar", Editar)  
	http.HandleFunc("/actualizar", Actualizar)
	log.Println("servidor corriendo...")         //iniciando el servidor
	log.Fatal(http.ListenAndServe(":8081", nil)) //si no se inicia el servidor, termina el programa
}
type Empleados struct { 
	Id	 int    
	Nombre string 
	Correo string 	
}
func Inicio(w http.ResponseWriter, r *http.Request) { //iniciando el servidor, w es la respuesta y r es la peticion
	// fmt.Fprintf(w, "Hola Mundo")
	//para leer la info necesitamos crear ESTRUCTURAS 
	
	conexionEstablecida := conexionBD()                                                                                              //llama a la funcion conexionBD
	registros, err := conexionEstablecida.Query("SELECT * FROM empleados") //preparar la consulta
	if err != nil {                                                                                                                  //si hay error
		panic(err) //muestra el error
	}
	empleado := Empleados{} //crea una variable de tipo empleados
	arrrgloEmpleado := []Empleados{} //crea un arreglo de tipo empleados

	for registros.Next() { //mientras haya registros
		var id int //crea una variable de tipo entero
		var nombre, correo string //crea una variable de tipo string	
		err := registros.Scan(&id, &nombre, &correo) //escanea los registros y los guarda en las variables
		if err != nil {                                                                                                                  //si hay error
			panic(err) //muestra el error
		}
		empleado.Id = id //asigna el id a la variable empleado
		empleado.Nombre = nombre //asigna el nombre a la variable empleado	
		empleado.Correo = correo //asigna el correo a la variable empleado
		arrrgloEmpleado = append(arrrgloEmpleado, empleado) //agrega el empleado al arreglo de empleados
	}
	
	plantillas.ExecuteTemplate(w, "inicio", arrrgloEmpleado) 

}

func Crear(w http.ResponseWriter, r *http.Request) { //meto similar a inicio
	plantillas.ExecuteTemplate(w, "crear", nil) //ruta de crear, muestra la plantilla crear

}

func Insertar(w http.ResponseWriter, r *http.Request) { //funcion para insertar datos en la base de datos
	if r.Method == "POST" { //si el metodo es post
		nombre := r.FormValue("nombre") //obtiene el valor del formulario
		correo := r.FormValue("correo") //obtiene el valor del formulario
		conexionEstablecida:= conexionBD()                                                                                              //llama a la funcion conexionBD
		insertarRegistro,err:= conexionEstablecida.Prepare("INSERT INTO empleados(nombre, correo) VALUES (?,?)") //preparar la consulta
		if err != nil { panic(err.Error()) }//lanzamos el error
			
		insertarRegistro.Exec(nombre, correo) //ejecuta la consultasql

		http.Redirect(w, r, "/", http.StatusMovedPermanently) //redirecciona a la pagina de inicio
	}
}

func Borrar(w http.ResponseWriter, r *http.Request) { //funcion para borrar datos de la base de datos
	idEmpleado := r.URL.Query().Get("id") //obtiene el id del empleado de la url de inicio.tmpl
	conexionEstablecida:= conexionBD()                                                                                              //llama a la funcion conexionBD
	borrarRegistro,err:= conexionEstablecida.Prepare("DELETE FROM empleados WHERE id=?") //preparar la consulta
	if err != nil { panic(err.Error()) }//lanzamos el error
			
	borrarRegistro.Exec(idEmpleado) //ejecuta la consultasql

	http.Redirect(w, r, "/", http.StatusMovedPermanently) //redirecciona a la pagina de inicio
}

func Editar(w http.ResponseWriter, r *http.Request) { //funcion para editar datos de la base de datos
	idEmpleado := r.URL.Query().Get("id") //obtiene el id del empleado de la url de inicio.tmpl
	conexionEstablecida:= conexionBD() 
	empleado := Empleados{} //crea una variable de tipo empleados 
	registro, err := conexionEstablecida.Query("SELECT * FROM empleados WHERE id=?", idEmpleado) //preparar la consulta 
	if err != nil {
		panic(err) //muestra el error
	}
	defer registro.Close() //asegura que se cierre el registro después de usarlo

	for registro.Next() { //mientras haya registros
		var id int //crea una variable de tipo entero
		var nombre, correo string //crea una variable de tipo string	
		err := registro.Scan(&id, &nombre, &correo) //escanea los registros y los guarda en las variables
		if err != nil {                                                                                                                  
			panic(err) //muestra el error
		}
		empleado.Id = id //asigna el id a la variable empleado
		empleado.Nombre = nombre //asigna el nombre a la variable empleado	
		empleado.Correo = correo //asigna el correo a la variable empleado
	}
	fmt.Println(empleado) //muestra el empleado
	plantillas.ExecuteTemplate(w, "editar", empleado) //muestra la plantilla editar

}
func Actualizar(w http.ResponseWriter, r *http.Request) { //funcion para insertar datos en la base de datos
	if r.Method == "POST" { //si el metodo es post
		id := r.FormValue("id") 
		nombre := r.FormValue("nombre") //obtiene el valor del formulario
		correo := r.FormValue("correo") //obtiene el valor del formulario
		conexionEstablecida:= conexionBD()                                                                                              //llama a la funcion conexionBD
		modificarRegistro,err:= conexionEstablecida.Prepare(" UPDATE empleados SET nombre=?, correo=? WHERE id=?") //preparar la consulta
		if err != nil { panic(err.Error()) }//lanzamos el error
			
		modificarRegistro.Exec(nombre, correo, id) //ejecuta la consultasql

		http.Redirect(w, r, "/", http.StatusSeeOther) //redirecciona a la pagina de inicio
	}
}