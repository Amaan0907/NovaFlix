import React from 'react'
import Movie from '../movie/Movie'

const Movies = ({ movies, message }) => {
    return (
        <div className='container mt-4 '>
            <div className='row'>
                {movies && movies.length > 0 ? movies.map((movie,id) => (
                    <Movie key={id} movie={movie} />
                )): <h2>{message}</h2>
}
            </div>

        </div>
    )
}

export default Movies
