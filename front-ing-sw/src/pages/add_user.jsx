import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import { toast, ToastContainer } from "react-toastify";
import "react-toastify/dist/ReactToastify.css";
import Navbar from "../components/navbar";

const AddUser = () => {
    const [user, setUser] = useState({
        name: "",
        last_name: "",
        username: "",
        password: "",
        email: "",
        phone: "",
        address: ""
    });
    
    const navigate = useNavigate();

    const handleChange = (e) => {
        const { name, value } = e.target;

        // Convert phone to an integer only if the input is for the phone field
        if (name === "phone") {
            const intPhone = Number(value); // Remove non-numeric chars and convert to int
            setUser({ ...user, [name]: intPhone || '' }); // Store as an integer if valid, else empty string
        } else {
            setUser({ ...user, [name]: value });
        }
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        const baseUrl = process.env.REACT_APP_API_BASE_URL;

        // Validate the form before submission
        if (!validateForm()) {
            return;
        }

        try {
            const request = await fetch(`${baseUrl}/user`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify(user)
            });

            const response = await request.json();

            console.log(response)

            if (!request.ok) {

                toast.error(`${response.message}`);
                
            } else {

                toast.success("Usuario agregado con éxito");
                setTimeout(() => {
                    navigate("/")
                }, 3000); // Delays for 3 seconds

            }

        } catch (error) {
            toast.error(`Error: ${error}`);
        }
    };

    const validateForm = () => {
        const { name, last_name, username, email, phone, address, password } = user;
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/; // Simple email regex
        const passwordRegex = /^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d]{8,}$/; // Password should have at least 8 characters, one letter, one number

        if (!name || !last_name || !username || !email || !phone || !address || !password) {
            toast.error("Por favor, complete todos los campos.");
            return false;
        }
        if (!emailRegex.test(email)) {
            toast.error("Por favor, ingrese un correo electrónico válido.");
            return false;
        }
        if (!passwordRegex.test(password)) {
            toast.error("La contraseña debe tener al menos 8 caracteres, incluyendo una letra y un número.");
            return false;
        }
        return true;
    };

    return (
        <div className="App">
            <Navbar />
            <div className="add-user-container">
            <h2>Agregar Nuevo Usuario</h2>
            <form onSubmit={handleSubmit} className="add-user-form">
                <div className="form-group">
                    <label htmlFor="name">Nombre:</label>
                    <input
                        id="name"
                        type="text"
                        name="name"
                        value={user.name}
                        onChange={handleChange}
                        required
                    />
                </div>
                <div className="form-group">
                    <label htmlFor="last_name">Apellido:</label>
                    <input
                        id="last_name"
                        type="text"
                        name="last_name"
                        value={user.last_name}
                        onChange={handleChange}
                        required
                    />
                </div>
                <div className="form-group">
                    <label htmlFor="username">Username:</label>
                    <input
                        id="username"
                        type="text"
                        name="username"
                        value={user.username}
                        onChange={handleChange}
                        required
                    />
                </div>
                <div className="form-group">
                    <label htmlFor="password">Contraseña (debe tener al menos 1 numero y 8 caracteres):</label>
                    <input
                        id="password"
                        type="password"
                        name="password"
                        value={user.password}
                        onChange={handleChange}
                        required
                    />
                </div>
                <div className="form-group">
                    <label htmlFor="email">Email:</label>
                    <input
                        id="email"
                        type="email"
                        name="email"
                        value={user.email}
                        onChange={handleChange}
                        required
                    />
                </div>
                <div className="form-group">
                    <label htmlFor="phone">Teléfono (debe tener solo 7 caracteres. No incluir codigo de area ni de pais):</label>
                    <input
                        id="phone"
                        type="text"
                        name="phone"
                        value={user.phone}
                        onChange={handleChange}
                        required
                    />
                </div>
                <div className="form-group">
                    <label htmlFor="address">Dirección:</label>
                    <input
                        id="address"
                        type="text"
                        name="address"
                        value={user.address}
                        onChange={handleChange}
                        required
                    />
                </div>
                <button type="submit" className="submit-button">Agregar Usuario</button>
            </form>
            </div>
            <ToastContainer />
        </div>
    );
};

export default AddUser;