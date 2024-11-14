import React, { useState, useEffect } from "react";
import Navbar from "../components/navbar";
import { useNavigate } from 'react-router-dom';
import { ToastContainer, toast } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';

const HomePage = () => {
    const [users, setUsers] = useState([]);
    const [searchTerm, setSearchTerm] = useState('');

    const navigate = useNavigate();

    const getUsers = async () => {
        const baseUrl = process.env.REACT_APP_API_BASE_URL;

        try {
            const request = await fetch(`${baseUrl}/user`);
            const response = await request.json();
            setUsers(response);
          } catch (error) {
            console.log("No se pudieron obtener los usuarios:", error);
          }

    }

    useEffect(() => {
        getUsers();
      }, []);

    // Filter users based on the search term
    const filteredUsers = users.filter(user => 
        user.name.toLowerCase().includes(searchTerm.toLowerCase())
    );

    // Handle search input change
    const handleSearchChange = (event) => {
        setSearchTerm(event.target.value);
    };

    const handleModify = (id) => {
        console.log('Navigating to:', `/modify-user/${id}`);
        navigate(`/modify-user/${id}`);
    }

    const handleEliminate = async (id) => {
        const baseUrl = process.env.REACT_APP_API_BASE_URL;

         // Show confirmation dialog
         const confirmDelete = window.confirm("Estas seguro que quieres eliminar este usuario?");
        
         if (!confirmDelete) {
             return; // If user clicked "Cancel", exit the function
         }

        try {
            const response = await fetch(`${baseUrl}/user/${id}`, {
                method: 'DELETE', // Specify the HTTP method as DELETE
            });
    
            if (!response.ok) {
                // Check if the response is not okay (status code outside 200-299)
                throw new Error(`Error: ${response.statusText}`);
            }

            toast.success("Usuario eliminado con éxito");
                setTimeout(() => {
                    // Optionally, refresh the user list after deletion
                    getUsers(); // Fetch the updated user list after deletion
                }, 1500); 

    
          } catch (error) {
            console.log("Error al eliminar usuario", error);
            toast.error("Error al eliminar usuario", error)
          }
    }

    return(
        <div className="App">
            <Navbar />
            <div className="search-container">
                <input
                type="text"
                className="search-bar"
                placeholder="Busca por nombre..."
                value={searchTerm}
                onChange={handleSearchChange}
                />
            </div>
            <div className="users-container">
            {filteredUsers.length ? (
                filteredUsers.map((user) => (
                <div className='user-card' key={user.id}>
                    <div className="user-info">
                    <h4>{user.name} {user.last_name}</h4>
                    <p>Username: {user.username}</p>
                    <p>Email: {user.email}</p>
                    <p>Telefono: {user.phone}</p>
                    <p>Direccion: {user.address}</p>
                    </div>
                    <button className="user-button" onClick={() => handleModify(user.id)}>Modificar</button>
                    <button className="user-button" onClick={() => handleEliminate(user.id)}>Eliminar</button>
                </div>
                ))
            ) : (
                <p>No hay usuarios</p>
            )}
            </div>
            <ToastContainer />
        </div>
    );

};

export default HomePage;