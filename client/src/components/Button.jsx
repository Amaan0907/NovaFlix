function Button({ children, onClick, type = "button", className = "" }) {
  return (
    <button
      type={type}
      onClick={onClick}
      className={`px-4 py-2 rounded-md bg-red-600 text-white font-medium hover:bg-red-700 transition-colors ${className}`}
    >
      {children}
    </button>
  );
}

export default Button;
