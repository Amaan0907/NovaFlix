import { Link } from "react-router-dom";
import Container from "./Container";

function Navbar({ children }) {
  return (
    <nav className="bg-black text-white shadow-2xs">
      <Container className="flex items-center justify-between py-3">
        <div className="flex items-center gap-8">
          <Link to="/" className="text-2xl font-bold text-red-600">NovaFlix</Link>

          <ul className="flex items-center gap-6">
            <li><Link to="/" className="hover:text-gray-300">Home</Link></li>
            <li><Link to="/movies" className="hover:text-gray-300">Movies</Link></li>
            <li><Link to="/tv-shows" className="hover:text-gray-300">TV Shows</Link></li>
            <li><Link to="/my-list" className="hover:text-gray-300">My List</Link></li>
          </ul>
        </div>

        {children}
      </Container>
    </nav>
  );
}

export default Navbar;
