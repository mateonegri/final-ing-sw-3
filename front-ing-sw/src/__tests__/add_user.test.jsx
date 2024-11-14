import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import AddUser from '../pages/add_user';
import { toast } from 'react-toastify';

// Mock the useNavigate function from react-router-dom
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: jest.fn(() => jest.fn()),
}));

// Mock the toast notifications
jest.mock('react-toastify', () => ({
  toast: {
    success: jest.fn(),
    error: jest.fn(),
  },
  ToastContainer: () => <div />,
}));

describe('AddUser component', () => {
  beforeEach(() => {
    global.fetch = jest.fn(() =>
      Promise.resolve({
        ok: true,
        json: () => Promise.resolve({ message: 'User added successfully' }),
      })
    );
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

  test('renders the AddUser form', () => {
    render(
      <MemoryRouter>
        <AddUser />
      </MemoryRouter>
    );

    expect(screen.getByText('Agregar Nuevo Usuario')).toBeInTheDocument();
    expect(screen.getByLabelText('Nombre:')).toBeInTheDocument();
    expect(screen.getByLabelText('Apellido:')).toBeInTheDocument();
  });

  test('handles input changes', () => {
    render(
      <MemoryRouter>
        <AddUser />
      </MemoryRouter>
    );

    const nameInput = screen.getByLabelText('Nombre:');
    fireEvent.change(nameInput, { target: { value: 'John', name: 'name' } });
    expect(nameInput.value).toBe('John');
  });

  test('shows an error when form is incomplete', async () => {
    render(
      <MemoryRouter>
        <AddUser />
      </MemoryRouter>
    );

    const submitButton = screen.getByRole('button', { name: /agregar/i });
    fireEvent.click(submitButton);

    await waitFor(() => {
      expect(toast.error).toHaveBeenCalledWith('Por favor, complete todos los campos.');
    });
  });

  test('submits the form with valid data', async () => {
    render(
      <MemoryRouter>
        <AddUser />
      </MemoryRouter>
    );

    fireEvent.change(screen.getByLabelText('Nombre:'), { target: { value: 'John', name: 'name' } });
    fireEvent.change(screen.getByLabelText('Apellido:'), { target: { value: 'Doe', name: 'last_name' } });
    fireEvent.change(screen.getByLabelText('Username:'), { target: { value: 'johndoe', name: 'username' } });
    fireEvent.change(screen.getByLabelText('Contraseña (debe tener al menos 1 numero y 8 caracteres):'), {
      target: { value: 'password1', name: 'password' },
    });
    fireEvent.change(screen.getByLabelText('Email:'), { target: { value: 'john@example.com', name: 'email' } });
    fireEvent.change(screen.getByLabelText('Teléfono (debe tener solo 7 caracteres. No incluir codigo de area ni de pais):'), {
      target: { value: '1234567', name: 'phone' },
    });
    fireEvent.change(screen.getByLabelText('Dirección:'), { target: { value: '123 Street', name: 'address' } });

    fireEvent.click(screen.getByRole('button', { name: /agregar/i }));

    await waitFor(() => {
      expect(global.fetch).toHaveBeenCalledWith(`${process.env.REACT_APP_API_BASE_URL}/user`, expect.any(Object));
      expect(toast.success).toHaveBeenCalledWith('Usuario agregado con éxito');
    });
  });
});
