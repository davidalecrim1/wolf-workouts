export const NavBar = () => {
    return (
        <nav className="w-full bg-white shadow-md">
            <div className="px-6 py-4">
                <div className="flex justify-between items-center">
                    <span className="text-xl">Wolf Workouts </span>
                    <div className="flex items-center gap-4">
                        <div className="flex items-center gap-4">
                            <a href="/trainings">Trainings</a>
                            <a href="/calendar">Calendar</a>
                            <a href="/schedule">Set Schedule</a>
                            <button className="border-2 border-blue-500 text-blue-500 hover:bg-blue-500 hover:text-white px-4 py-2 rounded">Logout</button>
                        </div>
                    </div>
                </div>
            </div>
        </nav>
    )
}
