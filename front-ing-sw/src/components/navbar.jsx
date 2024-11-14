import React from "react";
import "../styles/navbar.css";

const Navbar = () => {

    return (

        <nav className="navbar">
            <div className="navbar-logo">
                <a href="/">
                    Usuarios
                </a>
            </div>
            <ul className="navbar-links">
                <li><a href="/">Inicio</a></li>
                <li><a href="/add-user">Agregar Usuario</a></li>
            </ul>
        </nav>
    );
};

export default Navbar;