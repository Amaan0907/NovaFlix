function Container({ children, fluid, className = "" }) {
  return (
    <div
      className={
        fluid
          ? `w-full mx-auto px-4 ${className}`
          : `mx-auto px-4 sm:max-w-135 md:max-w-180 lg:max-w-240 xl:max-w-285 ${className}`
      }
    >
      {children}
    </div>
  );
}

export default Container;
