import React, { useState, useEffect } from "react";
import { useParams, useNavigate } from "react-router-dom";
import Navbar from "../components/navbar";
import { ToastContainer, toast } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';

const ModifyUser = () => {
    const { id } = useParams(); // Get the user ID from the URL parameters
    const [user, setUser] = useState({
        name: '',
        lastName: '',
        username: '',
        email: '',
        phone: '',
        address: ''
    });
    const [loading, setLoading] = useState(true); // State to handle loading status
    const [error, setError] = useState(null); // State to handle error messages
    const navigate = useNavigate();

    // Function to fetch user data from the API
    const getUser = async () => {
        const baseUrl = process.env.REACT_APP_API_BASE_URL;
        try {
            const response = await fetch(`${baseUrl}/user/${id}`);

            if (!response.ok) {

                toast.error(`Error loading user: ${error.message}`);

            } else {

                const data = await response.json();
                setUser(data);

            }

        } catch (err) {
            toast.error(`Error loading user: ${err.message || 'Unknown error'}`);
        } finally {
            setLoading(false);
        }
    };

    // Fetch user data when component mounts
    useEffect(() => {
        getUser();
    }, [id]);

    // Function to handle form submission
    const handleSubmit = async (event) => {
        event.preventDefault();
        const baseUrl = process.env.REACT_APP_API_BASE_URL;

        try {
            const request = await fetch(`${baseUrl}/user/${id}`, {
                method: "PUT",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(user), // Send the updated user data
            });

            const response = await request.json();

            console.log(response)

            if (request.ok) {
                toast.success("Usuario actualizado con éxito");
                setTimeout(() => {
                    navigate("/")
                }, 3000); // Delays for 3 seconds
            } else {
                toast.error(`${response.message}`);
            }

        } catch (err) {
            toast.error(`Error: ${error}`);
        }
    };

    // Handle input changes
    const handleChange = (event) => {
        const { name, value } = event.target;
        setUser((prevUser) => ({ ...prevUser, [name]: value }));
    };

    return (
        <div className="App">
            <Navbar />
                <div className="modify-user-container">
                <h2 className="modify-user-title">Modificar Usuario</h2>
                <form className="modify-user-form" onSubmit={handleSubmit}>
                    <div className="form-group">
                        <label htmlFor="name">Nombre:</label>
                        <input
                            id="name"
                            type="text"
                            name="name"
                            value={user.name || ''}
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
                            value={user.last_name || ''}
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
                            value={user.username || ''}
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
                            value={user.email || ''}
                            onChange={handleChange}
                            required
                        />
                    </div>
                    <div className="form-group">
                        <label htmlFor="phone">Telefono:</label>
                        <input
                            id="phone"
                            type="text"
                            name="phone"
                            value={user.phone || ''}
                            onChange={handleChange}
                            required
                        />
                    </div>
                    <div className="form-group">
                        <label htmlFor="address">Direccion:</label>
                        <input
                            id="address"
                            type="text"
                            name="address"
                            value={user.address || ''}
                            onChange={handleChange}
                            required
                        />
                    </div>
                    <button className="submit-button" type="submit" name="submit">Actualizar</button>
                </form>
                </div>
            <ToastContainer />
        </div>
    );
};

export default ModifyUser;
