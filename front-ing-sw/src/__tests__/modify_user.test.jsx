import { render, screen, fireEvent, waitFor, act } from '@testing-library/react';
import ModifyUser from '../pages/modify_user';
import '@testing-library/jest-dom';
import { ToastContainer, toast } from 'react-toastify';
import { MemoryRouter, Route, Routes } from 'react-router-dom';

// Mock the fetch API
global.fetch = jest.fn();

// Mock the toast notifications
jest.mock('react-toastify', () => ({
    toast: {
      success: jest.fn(),
      error: jest.fn(),
    },
    ToastContainer: () => <div />,
  }));

describe('ModifyUser Component', () => {
  beforeEach(() => {
    fetch.mockClear();
    jest.clearAllMocks();
  });

  test('renders ModifyUser form', async () => {
    
    await act(async () => {
        render(
          <MemoryRouter initialEntries={['/modify/1']}>
            <Routes>
              <Route path="/modify/:id" element={<ModifyUser />} />
            </Routes>
          </MemoryRouter>
        );
      });
      

    expect(screen.getByLabelText("Nombre:")).toBeInTheDocument();
    expect(screen.getByLabelText("Apellido:")).toBeInTheDocument();
    expect(screen.getByLabelText("Username:")).toBeInTheDocument();
    expect(screen.getByLabelText("Email:")).toBeInTheDocument();
    expect(screen.getByLabelText("Telefono:")).toBeInTheDocument();
    expect(screen.getByLabelText("Direccion:")).toBeInTheDocument();
});

  test('calls getUser on component mount', async () => {
    fetch.mockResolvedValueOnce({
      ok: true,
      json: async () => ({
        id: 1,
        name: 'John',
        last_name: 'Doe',
        username: 'johndoe',
        email: 'john@example.com',
        phone: '1234567890',
        address: '123 Main St',
      }),
    });

    render(
      <MemoryRouter initialEntries={['/modify-user/1']}>
        <Routes>
            <Route path="/modify-user/:id" element={<ModifyUser />} />
        </Routes>
      </MemoryRouter>
    );

    await waitFor(() => expect(fetch).toHaveBeenCalledTimes(1));
    expect(fetch).toHaveBeenCalledWith(expect.stringContaining('/user/1'));
  });

  test('displays user data after fetching', async () => {
    const mockUser = {
      id: 1,
      name: 'John',
      last_name: 'Doe',
      username: 'johndoe',
      email: 'john@example.com',
      phone: '1234567890',
      address: '123 Main St',
    };

    fetch.mockResolvedValueOnce({
      ok: true,
      json: async () => mockUser,
    });

    render(
      <MemoryRouter initialEntries={['/modify-user/1']}>
        <Routes>
            <Route path="/modify-user/:id" element={<ModifyUser />} />
        </Routes>
      </MemoryRouter>
    );

    await waitFor(() => {
        expect(fetch).toHaveBeenCalledTimes(1)
        expect(screen.getByDisplayValue(mockUser.name)).toBeInTheDocument();
        expect(screen.getByDisplayValue(mockUser.last_name)).toBeInTheDocument();
        expect(screen.getByDisplayValue(mockUser.username)).toBeInTheDocument();
        expect(screen.getByDisplayValue(mockUser.email)).toBeInTheDocument();
        expect(screen.getByDisplayValue(mockUser.phone)).toBeInTheDocument();
        expect(screen.getByDisplayValue(mockUser.address)).toBeInTheDocument();
    });
  });

  test('handles input changes', async () => {
    const mockUser = {
      name: 'John',
      last_name: 'Doe',
      username: 'johndoe',
      email: 'john@example.com',
      phone: '1234567890',
      address: '123 Main St',
    };

    fetch.mockResolvedValueOnce({
      ok: true,
      json: async () => mockUser,
    });

    render(
        <MemoryRouter initialEntries={['/modify-user/1']}>
            <Routes>
                <Route path="/modify-user/:id" element={<ModifyUser />} />
            </Routes>
        </MemoryRouter>
    );

    await waitFor(() => {
        expect(fetch).toHaveBeenCalledTimes(1)

        const nameInput = screen.getByLabelText(/nombre/i);
        fireEvent.change(nameInput, { target: { value: 'Jane' } });

        expect(nameInput.value).toBe('Jane');
    });
     
  });

  test('submits form successfully', async () => {
    const mockUser = {
      name: 'John',
      last_name: 'Doe',
      username: 'johndoe',
      email: 'john@example.com',
      phone: '1234567890',
      address: '123 Main St',
    };

    fetch.mockResolvedValueOnce({
      ok: true,
      json: async () => mockUser,
    }).mockResolvedValueOnce({
        ok: true,
        json: async () => mockUser,
      });

    render(
        <MemoryRouter initialEntries={['/modify-user/1']}>
            <Routes>
                <Route path="/modify-user/:id" element={<ModifyUser />} />
            </Routes>
        </MemoryRouter>
    );

    await waitFor(() => expect(fetch).toHaveBeenCalledTimes(1));

    const submitButton = screen.getByText("Actualizar")

    fireEvent.click(submitButton);

    await waitFor(() => {
        expect(fetch).toHaveBeenCalledTimes(2)
        expect(fetch).toHaveBeenCalledWith(
            expect.stringContaining('/user/1'),
            expect.objectContaining({
                method: 'PUT',
                body: JSON.stringify(mockUser),
        })
        );

        expect(toast.success).toHaveBeenCalledWith('Usuario actualizado con éxito');
    });
  });

  test('handles form submission error', async () => {
    fetch.mockRejectedValueOnce(new Error('Network error'));

    render(
        <MemoryRouter initialEntries={['/modify-user/1']}>
            <Routes>
                <Route path="/modify-user/:id" element={<ModifyUser />} />
            </Routes>
        </MemoryRouter>
    );

    await waitFor(() => expect(fetch).toHaveBeenCalledTimes(1));

    const submitButton = screen.getByRole('button', { name: /actualizar/i });
    fireEvent.click(submitButton);

    await waitFor(() => expect(toast.error).toHaveBeenCalledWith('Error: null'));
  });
});
